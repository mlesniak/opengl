#version 330 core
layout (location = 0) in vec3 pos;
layout (location = 1) in vec3 aNormal;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

uniform vec3 color;

out vec3 FragPos;
out vec3 objCol;
out vec3 Normal;

void main() {
    vec4 x;

    gl_Position = projection * view * model * vec4(pos, 1.0);
    objCol = color;
    x = model * vec4(aNormal, 0.0);
    Normal = vec3(x.x, x.y, x.z);
    FragPos = vec3(model * vec4(pos, 1.0));
}