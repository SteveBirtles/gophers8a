#version 430

out vec2 vertexCoord;

void main() {

    vec2 position = 2.0 * vec2((gl_VertexID % 2) * 2, (gl_VertexID / 2) * 2) - 1.0;
    gl_Position = vec4(position, 0.0, 1.0);
    vertexCoord = position;
}