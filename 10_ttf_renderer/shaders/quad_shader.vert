#version 450

layout(location=0) in vec2 inPosition;
layout(location=1) in vec3 inBary;

layout(location=0) out vec3 outBary;


void main() {
    gl_Position = vec4(inPosition[0]/320-0.8, inPosition[1]/320+0.8, 0.0, 1.0);
    outBary = inBary;
}