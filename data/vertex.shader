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
    // https://learnopengl.com/Lighting/Basic-Lighting
    // In the diffuse lighting section the lighting was fine because we didn't do any scaling
    // on the object, so there was not really a need to use a normal matrix and we could've
    // just multiplied the normals with the model matrix. If you are doing a non-uniform scale
    // however, it is essential that you multiply your normal vectors with the normal matrix.
    x = model * vec4(aNormal, 0.0);
    Normal = vec3(x.x, x.y, x.z);
    FragPos = vec3(model * vec4(pos, 1.0));
}