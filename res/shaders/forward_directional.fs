#version 410 core

uniform sampler2D diffuse;
uniform vec3 lightColor;

in vec2 TexCoord;
in vec3 LightPos;
in vec3 Normal;
in vec3 FragPos;

out vec4 fragColor;

float specularStrength = 0.5;

void main() {

    vec3 norm = normalize(Normal);

    vec3 lightDir = normalize(LightPos);

    float diff = max(dot(norm, lightDir), 0.0);

    vec3 diffuseLight = diff * lightColor;

    vec3 viewDir = normalize(-FragPos);
    vec3 halfwayDir = normalize(lightDir + viewDir);
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(norm, halfwayDir), 0.0), 128);
    vec3 specular = specularStrength * spec * lightColor;

    fragColor = texture(diffuse, TexCoord) * vec4(diffuseLight, 1.0f);
}
