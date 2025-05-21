// Basic OpenSocial JavaScript API mock/stub
// This will be expanded as needed.

window.gadgets = window.gadgets || {};

gadgets.util = gadgets.util || {};
gadgets.util.registerOnLoadHandler = function(callback) {
    if (document.readyState === "complete" || document.readyState === "interactive") {
        callback();
    } else {
        window.addEventListener('DOMContentLoaded', callback);
    }
};

gadgets.window = gadgets.window || {};
gadgets.window.adjustHeight = function(opt_height) {
    // In a real container, this would send a message to the parent to resize the iframe.
    // For now, log it. We might need a postMessage mechanism later.
    console.log("gadgets.window.adjustHeight called. Requested height:", opt_height || 'auto');
    // If opt_height is not provided, it means auto-adjust to content.
    // This is a simplified approach. A robust solution uses ResizeObserver or similar.
    let height = opt_height || document.body.scrollHeight;
    // This is a placeholder for cross-origin communication if domains differ
    if (window.parent && window.parent !== window) {
        // Example: window.parent.postMessage({ type: 'resize', height: height }, '*');
    }
};

// Placeholder for opensocial data APIs
window.opensocial = window.opensocial || {};
opensocial.newFetchPersonRequest = function(userId) {
    // This is a mock. It doesn't actually fetch.
    // It returns an object that mimics part of the OpenSocial API.
    return {
        userId: userId,
        __type: 'FetchPersonRequest' // Internal marker for our mock
    };
};

opensocial.newDataRequest = function() {
    let items = [];
    return {
        add: function(request, key) {
            items.push({ key: key, request: request });
        },
        send: function(callback) {
            // Simulate an async response
            setTimeout(function() {
                let responseData = {};
                items.forEach(function(item) {
                    if (item.request.__type === 'FetchPersonRequest') {
                        // Mocked response for FetchPersonRequest
                        responseData[item.key] = {
                            data: { // This structure mimics OpenSocial Person object
                                getId: function() { return item.request.userId; },
                                getDisplayName: function() { return "Mock " + item.request.userId; },
                                // Add other fields as needed by gadgets
                            },
                            // Helper for checking errors, common in OpenSocial API
                            hadError: function() { return false; } 
                        };
                    }
                });
                callback(responseData);
            }, 100);
        }
    };
};

console.log('Custom gadgets.js loaded for gadget.');
