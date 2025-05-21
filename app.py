from flask import Flask, render_template, abort, send_from_directory, request
import requests # Add this import for fetching external gadget URLs
from xml.etree import ElementTree as ET # For parsing gadget XML

app = Flask(__name__)

games_data = {
    "sample_game": {
        "name": "Sample Game",
        "description": "This is a sample game for demonstration purposes.",
        "developer": "Platform Team",
        "version": "1.0",
        "icon": "static/icons/sample_game_icon.png",
        "gadget_url": "/gadgets/sample_gadget.xml" # Updated to local path
    },
    "another_game": { # This one still points to an external example
        "name": "Another Exciting Game",
        "description": "A thrilling adventure awaits!",
        "developer": "Creative Studio",
        "version": "0.9beta",
        "icon": "static/icons/another_game_icon.png",
        "gadget_url": "http://www.google.com/ig/modules/googletalk.xml" # Example external gadget
    },
    "mini_game": {
        "name": "Mini Test Gadget",
        "description": "A very simple gadget for testing the container.",
        "developer": "Platform Team",
        "version": "0.1",
        "icon": None, # No icon for this one
        "gadget_url": "/gadgets/mini_gadget.xml"
    }
}

@app.route('/')
def top_page():
    return render_template('top.html', games=games_data)

@app.route('/game/<game_id>')
def game_intro_page(game_id):
    game = games_data.get(game_id)
    if not game:
        abort(404)
    return render_template('game_intro.html', game=game, game_id=game_id)

@app.route('/game/<game_id>/play')
def game_run_page(game_id):
    game = games_data.get(game_id)
    if not game:
        abort(404)
    return render_template('game_run.html', game=game, game_id=game_id)

@app.route('/gadgets/<path:filename>')
def serve_gadget(filename):
    # This serves gadget XML files from the game_data directory
    # In a production environment, you might want more robust access controls
    return send_from_directory('game_data', filename)

@app.route('/gadgets/ifr')
def gadget_iframe_content():
    gadget_url = request.args.get('url')
    if not gadget_url:
        abort(400, "Gadget URL not provided.")

    # Determine if the gadget_url is internal or external
    # For simplicity, assume internal gadgets start with '/'
    # and are served by our 'serve_gadget' route.
    # External URLs will be fetched using 'requests'.
    
    content_xml_str = ""
    try:
        if gadget_url.startswith('/'): # Internal gadget
            # Construct absolute URL for internal fetching
            # request.url_root gives something like 'http://127.0.0.1:5000/'
            internal_gadget_abs_url = request.url_root.rstrip('/') + gadget_url
            response = requests.get(internal_gadget_abs_url)
            response.raise_for_status() # Raise an exception for HTTP errors
            content_xml_str = response.text
        else: # External gadget (basic fetch)
            response = requests.get(gadget_url)
            response.raise_for_status()
            content_xml_str = response.text
    except requests.exceptions.RequestException as e:
        abort(500, f"Failed to fetch gadget XML: {e}")
    except Exception as e: # Catch other potential errors during fetch
        abort(500, f"An unexpected error occurred while fetching gadget XML: {e}")

    try:
        # Parse the gadget XML to extract the Content section
        # A more robust parser would handle namespaces, etc.
        root = ET.fromstring(content_xml_str)
        content_node = root.find(".//Content[@type='html']")
        if content_node is None or content_node.text is None:
            abort(500, "Gadget XML does not contain valid HTML content.")
        
        gadget_html_content = content_node.text
    except ET.ParseError as e:
        abort(500, f"Failed to parse gadget XML: {e}")
    except Exception as e: # Catch other potential errors during parsing
        abort(500, f"An unexpected error occurred while parsing gadget XML: {e}")


    # Render a template that includes gadgets.js and the gadget's HTML content
    return render_template('gadget_wrapper.html', gadget_html_content=gadget_html_content)

if __name__ == '__main__':
    app.run(debug=True)
