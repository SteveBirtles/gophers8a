#version 430

in vec2 vertexCoord;
out vec4 fragColor;

uniform int iterations;
uniform vec2 camera;
uniform float zoom;
uniform float power;
uniform float infinity;

vec3 mapToColor(in int i) {

    if (i == iterations) {
        return vec3(0, 0, 0);
    }

    float value = pow(float(i) / iterations, power) * 6.0;

    if (value < 1) {
        return vec3(value, 0, 1);
    } else if (value < 2) {
        value -= 1;
        return vec3(1, value, 1 - value);
    } else if (value < 3) {
        value -= 2;
        return vec3(1 - value, 1, 0);
    } else if (value < 4) {
        value -= 3;
        return vec3(0, 1, value);
    } else if (value < 5) {
        value -= 4;
        return vec3(0, 1 - value, 1);
    } else {
        value -= 5;
        return vec3(0, 0, value);
    }

}

void main() {

    vec2 z = vec2(0, 0);
    vec2 c = camera + zoom * vertexCoord;

    int i = 0;
    while (dot(z, z) < infinity && i < iterations) {
        z = vec2(z.x*z.x - z.y*z.y, 2.0*z.x*z.y) + c;
        i++;
    }

    vec3 color = mapToColor(i);
    fragColor = vec4(color, 1.0);
}