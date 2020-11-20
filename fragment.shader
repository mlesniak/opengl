#version 330 core
out vec4 color;

in vec3 objCol;
in vec3 Normal;
in vec3 FragPos;

uniform vec3 lightPos;

vec3 ambient = vec3(0.3, 0.3, 0.3);

vec3 norm = normalize(Normal);
vec3 lightDir = normalize(lightPos - FragPos);

float diff = max(dot(norm, lightDir), 0.0);
vec3 diffuse = diff * vec3(1, 1, 1);
vec3 result = (ambient + diffuse) * objCol;


void main() {


    color = vec4(result, 1.0);
}