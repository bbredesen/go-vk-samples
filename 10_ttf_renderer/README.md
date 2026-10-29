# go-vk-samples #10 - ttf-renderer

This mini-Vulkan app written in Go demonstrates usage of a stencil buffer to render text from a TrueType font file.

## Drawing the glyph

TTF glyphs are described with a series of Bezier curves drawn through two anchor points, plus a control point. The renderer does not directly convert the glyph outline to the output image, but rather writes first to a stencil that is then used to "knock out" the shape of the glyph from a solid rectangle.

There are three graphics pipelines bound across two rendering subpasses which build the stencil and then color through the stencil.

### Subpass 0 - building the stencil

The first subpass has no color attachment and only builds a stencil buffer across two rendering pipelines. It does not output a visible image. 

The glyph outlines in the TTF file are a list of closed line segments plus control points defining Bezier curves. Pipeline 0 builds the rough outline of the glyph from the line segments, ignoring the curve control points for now. The pipeline leverages TTFs clockwise winding rule and draws a triangle fan from one corner of the bounding box, incrementing the stencil for front facing triangles and decrementing for back facing triangles. 

Inspecting the stencil at this point would show a blocky, but recognizable glyph. 

Pipeline 1 then refines the curves of the glyph. Each curve segment of a TrueType glyph is defined as a quadratic Bezier curve. To expand the curves, we next draw triangles into the stencil using each segment's anchor points plus the control point. The fragment shader then tests each fragment against the quadratic curve to increment or decrement the stencil.

### Subpass 1 - coloring through the stencil

The second subpass actually renders against the swapchain image, attaching it as the color output. The single pipeline in this subpass draws a single quad covering the glyph's bounding box with a flat white fragment shader. The quad is compared against the stencil and only non-zero fragments in the stencil are passed through to the fragment shader, leaving just the image of the glyph.

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

* "320" is hard-coded in several places, notably the shaders. This is half of the em-width requested from sfnt. The sfnt library modifies the glyph bounds for better visuals based on the requested rendering size. Selecting a small render and scaling up creates distorted images.

I arbitrarily chose an em-width of 640 pixels per em, labeled `ppem` in the code) to minimize the effect.
