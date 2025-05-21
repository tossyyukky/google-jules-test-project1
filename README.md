# OpenSocial Game Platform PoC

This project is a Proof of Concept for a social game platform utilizing the OpenSocial API.

## Features

*   **TOP Page**: Displays a list of available games.
*   **Game Introduction Page**: Shows details for a selected game.
*   **Game Execution Page**: Hosts and runs OpenSocial gadgets (games).
*   **Basic OpenSocial Container**:
    *   Serves gadget XML files.
    *   Provides a minimal JavaScript environment (`gadgets.*` and `opensocial.*` mocks) to allow gadgets to run.
    *   Can render gadgets served locally or fetch from external URLs (with limitations due to mock API completeness).

## Project Structure

*   `app.py`: The main Flask application.
*   `static/`: Contains static assets (CSS, JavaScript, images).
    *   `static/js/gadgets.js`: Mock OpenSocial JavaScript API.
    *   `static/icons/`: Placeholder game icons.
*   `templates/`: HTML templates for the Flask application.
    *   `top.html`: Main page listing games.
    *   `game_intro.html`: Game details page.
    *   `game_run.html`: Page that hosts the game gadget iframe.
    *   `gadget_wrapper.html`: HTML shell that includes `gadgets.js` and the gadget content.
*   `game_data/`: Contains game gadget XML files and related data.
    *   `sample_gadget.xml`: A sample game gadget demonstrating basic features.
    *   `mini_gadget.xml`: A very simple gadget for testing.
*   `tests/`: Contains PyTest unit and integration tests.
    *   `test_app.py`: Tests for the Flask application routes and logic.
*   `requirements.txt`: Python dependencies.
*   `README.md`: This file.

## Setup and Running

1.  **Clone the repository** (if applicable).

2.  **Create a virtual environment** (recommended):
    ```bash
    python -m venv venv
    source venv/bin/activate  # On Windows: venv\Scripts\activate
    ```

3.  **Install dependencies**:
    ```bash
    pip install -r requirements.txt
    ```

4.  **Run the Flask application**:
    ```bash
    python app.py
    ```
    The application will typically be available at `http://127.0.0.1:5000/`.

5.  **Running Tests**:
    To run the automated tests, ensure PyTest is installed (it's in `requirements.txt`) and then run:
    ```bash
    pytest
    ```

## Implemented OpenSocial Features (Mocks)

*   **Gadget Rendering**: Gadgets of `type="html"` are supported. Their content is extracted and rendered within an iframe.
*   **`gadgets.js` API (Partial Mock)**:
    *   `gadgets.util.registerOnLoadHandler(callback)`: Executes the callback when the DOM is ready.
    *   `gadgets.window.adjustHeight()`: Logs to console (basic height adjustment not fully implemented across origins without `postMessage`).
*   **`opensocial.js` API (Partial Mock)**:
    *   `opensocial.newDataRequest()`: Creates a new data request object.
    *   `request.add(opensocial.newFetchPersonRequest('VIEWER'), 'viewer')`: Adds a request to fetch viewer data.
    *   `request.send(callback)`: Simulates sending the request and returns mock data.
    *   `Person.getDisplayName()`: Returns a mock display name for the viewer.

## Future Development Ideas

*   More complete OpenSocial API implementation (Persistence, Activities, Friends).
*   User authentication and management.
*   Database integration for games and user data.
*   More robust gadget rendering and security (e.g., using Caja or similar sandboxing).
*   Dynamic game list from a database.
*   Improved UI/UX.
