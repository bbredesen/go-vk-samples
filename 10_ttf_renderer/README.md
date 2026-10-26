# ttf-renderer

This mini-Vulkan app written in Go demonstrates usage of a stencil buffer to render text from a TrueType font file.

There are three graphics pipelines bound across two rendering subpasses. The first pipeline fills the body of the glyph
(and knocks out interior areas) by drawing a triangle fan across the vertices of the glyph. TrueType specifies a
clockwise winding rule. For example, on an 'O' character, there is an outer polygon with the vertices given in clockwise
order, and an inner polygon going counter-clockwise. 

The stencil test leverages this winding rule by incrementing the stencil value for front facing triangles, and decrementing for back
facing triangles. If you were to insepect the stencil after this phase, glyphs would be recognizable, but very blocky.

The second pipeline fills in the curves of the glyph. Triangles are drawn between each sequential pair of vertices,
along with the control point for that segment. The shader uses barymetric coordinates of each fragment to determine if that
fragment should be discarded or drawn in the stencil. (See quad_shader.frag)

In the second subpass, a color attachment is added and a square representing the bounds of the glyph is drawn and tested
against the stencil.

## Usage

Run `go generate` in the project root to compile the shaders before running. Requires the glsc compiler, bundled with the Vulkan SDK.

Pass a TTF font filepath with the `-font` flag, or set the character to render with `-char`.

## Known Issues

* The stencil is tested against a pair of triangles matching the glyph bounds provided by sfnt. There are several
  font/glyph combinations that are not getting bounds from sfnt, or getting bounds that are far too small. The net
  effect is just a blank screen. It seems to be more common with non-letter glyphs in display fonts. For example, see:
  * Elephant - &
  * Algerian - $
  
  I have not investigated if this is coming from the font file or a bug in sfnt

* "320" is hard-coded in several places, notably the shaders. This is half of the em-width requested from sfnt. (sfnt
  modifies the glyph bounds for better visuals based on the requested rendering size.) I arbitrarily chose an em-width of 640 pixels
  per em, labeled `ppem` in the code) to minimize the effect. The fix is to get the true curves from the TTF
  files and scaling proprtionately, instead of relying on sfnt.

## Next Steps

* Live rendering a string, not just a single glyph. This is a slipery slope into typography.
* Background rendering the full font (or a subset) to texture memory on the GPU, then being able to print text with a
  bunch of textured quads.
  * Also, each glyph could be generated as a mipmap. Take note of the ppem parameter passed to sfnt.
