#version 450

layout(location=0) in vec3 baryCoords;

layout(location=0) out vec4 outColor;

void main() {
    // Inside of the Bezier curve based on this formula: (s/2+t)^2 < t
    float s = baryCoords.y;
    float t = baryCoords.x;
    float comp = (s/2+t)*(s/2+t);

    if (comp > t) {
        // Vertex shader is sending full triangles composed of two anchors and their control points. If this fragment is
        // outside of the curve, then discard it. Comment this section out to see the "block" rendering of the full
        // triangle, instead of the glyph curves.
        discard;
    }

    outColor = vec4(1,1,1,1);
}