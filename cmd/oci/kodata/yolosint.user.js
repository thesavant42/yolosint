// ==UserScript==
// @name         yolosint dockerhub search
// @namespace    http://tampermonkey.net/
// @version      1.3
// @description  Adds a button to launch OCI Explorer from Docker Hub tags and image pages
// @author       thesavant2
// @match        https://hub.docker.com/*
// @icon         https://raw.githubusercontent.com/thesavant42/dagdotdev/refs/heads/main/docs/favicon-60x60.png
// @grant        none
// ==/UserScript==

(function() {
    'use strict';

    // === Configuration ===
    const OCI_BASE_URL = "http://localhost:8042/?image=";

    // Small SVG Icon (Globe/Layer style) for the button
    const ociIcon = `<img src="https://raw.githubusercontent.com/thesavant42/dagdotdev/refs/heads/main/docdork-32.png" alt="Docker Dorker"/>`;

    /**
     * Creates the styled button element
     */
    function createButton(url) {
        const btn = document.createElement('a');
        btn.href = url;
        btn.target = '_blank';
        btn.innerHTML = ociIcon;
        btn.title = "Inspect in Docker Dorker";
        // Flat button styling
        btn.style.cssText = `
            display: inline-flex;
            align-items: center;
            justify-content: center;
            margin-left: 8px;
            text-decoration: none;
            vertical-align: middle;
            opacity: 0.7;
            transition: opacity 0.2s ease-in-out;
        `;

        btn.onmouseover = () => {
            btn.style.opacity = '1';
        };
        btn.onmouseout = () => {
            btn.style.opacity = '0.7';
        };
        return btn;
    }

    /**
     * 1. Handle "Tags" List Page
     * Target: <a data-testid="navToImage">tagname</a>
     */
    function processTagsList() {
        // Find all tag links that we haven't processed yet
        const tagLinks = document.querySelectorAll('a[data-testid="navToImage"]:not([data-oci-processed])');

        tagLinks.forEach(link => {
            link.setAttribute('data-oci-processed', 'true');
            const tag = link.textContent.trim();

            // The href usually looks like: /layers/namespace/repo/tag/images/sha256-...
            // We can extract the namespace and repo from this to be safe
            const href = link.getAttribute('href');
            if (!href) return;

            const parts = href.split('/');
            // parts expected: ["", "layers", "namespace", "repo", "tag", ...]
            if (parts.length >= 5) {
                const namespace = parts[2];
                const repo = parts[3];
                // Construct: namespace/repo:tag
                const imageRef = `${namespace}/${repo}:${tag}`;
                const finalUrl = `${OCI_BASE_URL}${encodeURIComponent(imageRef)}`;

                const btn = createButton(finalUrl);
                // Append the button to the parent <p> or container so it sits next to the link
                if(link.parentNode) {
                    link.parentNode.appendChild(btn);
                }
            }
        });
    }

    /**
     * 2. Handle Specific "Layer/Image" Page
     * Target: <h1 class="... MuiTypography-h3 ...">namespace/repo:tag</h1>
     * We want to target the specific SHA digest if possible.
     */
    function processLayerPage() {
        // Only run if we are on a layers page
        if (!window.location.pathname.includes('/layers/')) return;

        // Find the main header. The class list is specific, but H1 is usually unique enough in this context.
        const header = document.querySelector('h1[class*="MuiTypography-h3"]:not([data-oci-processed])');

        if (header) {
            header.setAttribute('data-oci-processed', 'true');

            // Parse current URL to get the specific Digest
            // URL format: /layers/namespace/repo/tag/images/sha256-DIGEST
            const pathParts = window.location.pathname.split('/');

            // Find the part that starts with sha256-
            const digestPart = pathParts.find(p => p.startsWith('sha256-'));

            if (digestPart && pathParts.length >= 4) {
                const namespace = pathParts[2];
                const repo = pathParts[3];
                // Convert "sha256-1234..." to "sha256:1234..."
                const digest = digestPart.replace('sha256-', 'sha256:');

                // Construct: namespace/repo@sha256:digest
                const imageRef = `${namespace}/${repo}@${digest}`;
                const finalUrl = `${OCI_BASE_URL}${encodeURIComponent(imageRef)}`; // Pass optional params if needed

                const btn = createButton(finalUrl);
                btn.style.marginLeft = '12px';
                btn.style.height = '24px'; // Slightly larger for the header
                btn.style.width = '24px';

                header.appendChild(btn);
            }
        }
    }

    // === Observer to handle SPA Navigation ===
    // Docker Hub is a Single Page App, so we need to watch for DOM changes
    const observer = new MutationObserver((mutations) => {
        // Run checks on DOM change
        processTagsList();
        processLayerPage();
    });

    observer.observe(document.body, {
        childList: true,
        subtree: true
    });

})();