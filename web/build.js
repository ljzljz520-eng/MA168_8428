const fs = require('fs');
const path = require('path');
const source = path.join(__dirname, 'index.html');
const outputDir = path.join(__dirname, 'dist');
fs.mkdirSync(outputDir, { recursive: true });
fs.copyFileSync(source, path.join(outputDir, 'index.html'));
fs.writeFileSync(path.join(outputDir, 'build-manifest.json'), JSON.stringify({name: 'bookstore-recommendation-web', version: '1.0.0'}));
