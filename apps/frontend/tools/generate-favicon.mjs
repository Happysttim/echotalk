// Build binary favicon assets from the editable SVG. No runtime dependency.
// Usage: node tools/generate-favicon.mjs [path-to-sharp-package]
import { createRequire } from 'node:module';
import { readFile, writeFile } from 'node:fs/promises';

const require = createRequire(import.meta.url);
const sharp = require(process.argv[2] || 'sharp');
const source = await readFile(
  new URL('../public/favicon.svg', import.meta.url),
);
const sizes = [16, 32, 48, 64, 256];
const frames = await Promise.all(
  sizes.map((size) =>
    sharp(source, { density: 384 }).resize(size, size).png().toBuffer(),
  ),
);
// ICO directory followed by PNG-compressed image entries (32-bit RGBA).
const directory = Buffer.alloc(6 + sizes.length * 16);
directory.writeUInt16LE(1, 2);
directory.writeUInt16LE(sizes.length, 4);
let offset = directory.length;
frames.forEach((png, index) => {
  const entry = 6 + index * 16;
  directory[entry] = sizes[index] === 256 ? 0 : sizes[index];
  directory[entry + 1] = directory[entry];
  directory.writeUInt16LE(1, entry + 4);
  directory.writeUInt16LE(32, entry + 6);
  directory.writeUInt32LE(png.length, entry + 8);
  directory.writeUInt32LE(offset, entry + 12);
  offset += png.length;
});
await writeFile(
  new URL('../public/favicon.ico', import.meta.url),
  Buffer.concat([directory, ...frames]),
);
await writeFile(
  new URL('../public/favicon.png', import.meta.url),
  frames.at(-1),
);
console.log(`Generated favicon.ico (${sizes.join(', ')}px) and favicon.png`);
