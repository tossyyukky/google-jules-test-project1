import pytest
from app import app as flask_app # Assuming your Flask app instance is named 'app' in app.py
import requests # Import for mocking

@pytest.fixture
def app():
    yield flask_app

@pytest.fixture
def client(app):
    return app.test_client()

def test_top_page(client):
    """Test the top page loads."""
    response = client.get('/')
    assert response.status_code == 200
    assert b"Available Games" in response.data
    assert b"Sample Game" in response.data # Check if sample_game is listed
    assert b"Mini Test Gadget" in response.data # Check if mini_game is listed

def test_game_intro_page_valid(client):
    """Test the game introduction page for a valid game."""
    response = client.get('/game/sample_game')
    assert response.status_code == 200
    assert b"Sample Game" in response.data
    assert b"Description:" in response.data
    assert b"Start Game" in response.data

def test_game_intro_page_invalid(client):
    """Test the game introduction page for an invalid game."""
    response = client.get('/game/non_existent_game')
    assert response.status_code == 404

def test_game_run_page_valid(client):
    """Test the game run page for a valid game."""
    response = client.get('/game/sample_game/play')
    assert response.status_code == 200
    assert b"Now Playing: Sample Game" in response.data
    # Check if the iframe source is correctly pointing to the gadget container
    assert b'src="/gadgets/ifr?url=%2Fgadgets%2Fsample_gadget.xml"' in response.data

def test_game_run_page_invalid(client):
    """Test the game run page for an invalid game."""
    response = client.get('/game/non_existent_game/play')
    assert response.status_code == 404

def test_serve_gadget_xml_valid(client):
    """Test serving a valid gadget XML file."""
    response = client.get('/gadgets/sample_gadget.xml')
    assert response.status_code == 200
    # Flask's send_from_directory usually infers Content-Type.
    # For .xml, it might be application/xml or text/xml.
    # Let's check if 'xml' is in the mimetype.
    assert 'xml' in response.mimetype.lower()
    assert b"<ModulePrefs title=\"Sample Game Gadget\"" in response.data

def test_serve_gadget_xml_invalid(client):
    """Test serving a non-existent gadget XML file."""
    response = client.get('/gadgets/non_existent_gadget.xml')
    assert response.status_code == 404 # send_from_directory returns 404

def test_gadget_ifr_route_valid_local(client):
    """Test the /gadgets/ifr route with a valid local gadget URL."""
    response = client.get('/gadgets/ifr?url=/gadgets/sample_gadget.xml')
    assert response.status_code == 200
    assert b"<title>Gadget Container</title>" in response.data
    assert b'<script src="/static/js/gadgets.js"></script>' in response.data
    # Check for some content from sample_gadget.xml
    assert b"Hello from Sample Game!" in response.data

def test_gadget_ifr_route_missing_url(client):
    """Test the /gadgets/ifr route with no URL parameter."""
    response = client.get('/gadgets/ifr')
    assert response.status_code == 400 # Bad Request

def test_gadget_ifr_route_fetch_error(client, monkeypatch):
    """Test the /gadgets/ifr route when fetching an external URL fails."""
    # Mock requests.get to simulate a fetch error
    class MockResponse:
        def __init__(self, status_code):
            self.status_code = status_code
            self.text = "Error"
        def raise_for_status(self):
            if self.status_code != 200:
                raise requests.exceptions.HTTPError("Simulated HTTP error")

    def mock_get(url, **kwargs):
        if "error_example.com" in url:
            raise requests.exceptions.RequestException("Simulated connection error")
        # For internal calls, we need to simulate the actual app's behavior
        # This part is tricky because the app itself uses requests.get for internal URLs.
        # A better mock would distinguish between external and internal, or we rely on other tests for internal.
        # For now, let's assume this test focuses on external fetch errors.
        # The test as written might fail if an internal URL is hit by this mock_get.
        # A more robust setup might involve patching at a lower level or using a more specific target for patching.
        
        # To avoid breaking internal calls made by the app to itself (e.g. /gadgets/sample_gadget.xml)
        # we need to ensure they don't raise the "Simulated connection error".
        # The original test code did not account for this.
        # A simple way for *this specific test case* is to let other URLs pass through unmocked,
        # but monkeypatch doesn't work that way directly.
        # The most direct fix is to ensure the app code correctly uses the patched 'requests.get'.
        # The provided test code for 'mock_get' was:
        # return MockResponse(200) # Default good response for other calls (like internal)
        # This is problematic because internal calls might not expect a MockResponse object.
        # The app's internal requests.get call expects a real Response object from the local server.
        #
        # For this test, we are only concerned with testing the *external* fetch error.
        # The app code does:
        #   internal_gadget_abs_url = request.url_root.rstrip('/') + gadget_url
        #   response = requests.get(internal_gadget_abs_url)
        # This means our mock_get WILL intercept this call.
        #
        # Given the constraints, the simplest way to make *this test pass* and focus on external errors,
        # is to make the mock_get specific to the erroring URL and let others pass (which is hard with setattr)
        # or make it return a valid-enough response for internal calls.
        # Let's assume the provided test logic intended to only mock for "error_example.com".
        # The original `requests.get` should be called for other URLs if `monkeypatch.setattr` is not specific enough.
        # However, `monkeypatch.setattr(requests, "get", mock_get)` replaces `requests.get` globally.
        #
        # Re-evaluating the mock_get:
        # If url contains "error_example.com", raise error.
        # Otherwise, it must be an internal call like "http://127.0.0.1:5000/gadgets/sample_gadget.xml".
        # The test for `/gadgets/ifr?url=/gadgets/sample_gadget.xml` (test_gadget_ifr_route_valid_local)
        # already covers the case where internal `requests.get` works.
        # This specific test `test_gadget_ifr_route_fetch_error` is for when `url` is external and fails.
        # So, the mock_get only needs to correctly handle the erroring external URL.
        # The `app.py` code:
        #   if gadget_url.startswith('/'): response = requests.get(internal_gadget_abs_url)
        #   else: response = requests.get(gadget_url)
        # So if `gadget_url` is `http://error_example.com/gadget.xml`, `requests.get` is called with this URL.
        # The `mock_get` will then correctly raise `RequestException`.
        # The provided `mock_get` seems fine for this specific test's purpose.
        raise requests.exceptions.RequestException("Simulated connection error")


    # Patch requests.get for the duration of this test
    monkeypatch.setattr(requests, "get", mock_get)
    
    response = client.get('/gadgets/ifr?url=http://error_example.com/gadget.xml')
    assert response.status_code == 500
    assert b"Failed to fetch gadget XML" in response.data
    assert b"Simulated connection error" in response.data # Check for our specific error

def test_gadget_ifr_route_bad_xml(client, monkeypatch):
    """Test the /gadgets/ifr route with bad XML content."""
    def mock_get_bad_xml(url, **kwargs):
        class MockResponse:
            def __init__(self):
                self.text = "<Module><Content type='html'>This is not well-formed" # Bad XML
                self.status_code = 200
            def raise_for_status(self):
                pass # Simulate no HTTP error
        return MockResponse()

    monkeypatch.setattr(requests, "get", mock_get_bad_xml)
    
    # The URL here needs to be one that would trigger the parsing logic.
    # It can be a "local" gadget URL path because the actual fetching is mocked.
    response = client.get('/gadgets/ifr?url=/gadgets/bad_xml_gadget_for_test.xml') 
    assert response.status_code == 500
    assert b"Failed to parse gadget XML" in response.data

# To run these tests, navigate to the project root in the terminal
# and run the command: pytest
