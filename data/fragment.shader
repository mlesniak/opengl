#version 330 core
out vec4 color;

in vec4 objCol;
in vec3 Normal;
in vec3 FragPos;

uniform vec3 lightPos;

vec4 ambient = vec4(0.3, 0.3, 0.3, 1.0);

vec3 norm = normalize(Normal);
vec3 lightDir = normalize(lightPos - FragPos);

float diff = max(dot(norm, lightDir), 0.0);
vec4 diffuse = diff * vec4(1, 1, 1, 1);
vec4 result = (ambient + diffuse) * objCol;


void main() {


    color = result;
}