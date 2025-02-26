// ==UserScript==
// @name         VirusTotal Subdomain Exporter
// @namespace    http://tampermonkey.net/
// @version      1.1
// @description  Export all subdomains from VirusTotal in JSON format
// @author       savant42
// @match        https://www.virustotal.com/gui/domain/*/relations
// @grant        GM_xmlhttpRequest
// @grant        GM_download
// @run-at       document-idle
// ==/UserScript==

(function () {
    'use strict';

    const API_KEY = "YOURKEYHERE"; // 🔑 Replace with your valid API key
    let domain = window.location.pathname.split('/')[3]; // Extract domain from URL

    async function fetchAllSubdomains(domain) {
        let url = `https://www.virustotal.com/api/v3/domains/${domain}/subdomains?limit=40`;
        let allResults = [];

        try {
            while (url) {
                console.log(`🚀 Fetching: ${url}`);
                let response = await fetch(url, {
                    headers: { 'x-apikey': API_KEY }
                });

                if (!response.ok) {
                    console.error(`❌ Error fetching data: ${response.status} ${response.statusText}`);
                    alert("Error fetching subdomains! Check API key & permissions.");
                    return null;
                }

                let data = await response.json();
                allResults.push(...data.data);

                url = data.links?.next || null; // Pagination handling
            }

            console.log(`✅ Successfully fetched ${allResults.length} subdomains!`);
            return allResults;
        } catch (error) {
            console.error("❌ Error fetching subdomains:", error);
            return null;
        }
    }

    function downloadJSON(data, filename) {
        const jsonStr = JSON.stringify(data, null, 2);
        const blob = new Blob([jsonStr], { type: 'application/json' });
        const url = URL.createObjectURL(blob);

        const a = document.createElement('a');
        a.href = url;
        a.download = filename;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);

        URL.revokeObjectURL(url);
    }

    function addDownloadButton() {
        setTimeout(() => {
            let toolbar = document.querySelector('.vt-ui-container') || document.body;
            let btn = document.createElement('button');

            btn.innerText = "Download Subdomains JSON";
            btn.style = "position: fixed; top: 10px; right: 10px; z-index: 9999; padding: 10px; background-color: #4CAF50; color: white; border: none; cursor: pointer;";
            btn.onclick = async () => {
                btn.innerText = "Fetching...";
                let subdomains = await fetchAllSubdomains(domain);
                if (subdomains) {
                    downloadJSON(subdomains, `${domain}_subdomains.json`);
                    btn.innerText = "Download Complete!";
                    setTimeout(() => (btn.innerText = "Download Subdomains JSON"), 3000);
                } else {
                    btn.innerText = "Error!";
                }
            };

            toolbar.appendChild(btn);
            console.log("✅ Download button added!");
        }, 2000);
    }

    addDownloadButton();
})();
