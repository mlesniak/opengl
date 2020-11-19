#version 330 core
out vec4 color;

in vec3 objCol;

void main() {
    color = vec4(objCol, 1.0);
}