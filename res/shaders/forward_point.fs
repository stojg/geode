#version 410

uniform sampler2D diffuse;
uniform vec3 lightPos;
uniform vec3 lightColor;
uniform vec3 viewPos;

in vec2 TexCoord;
in vec3 Normal;
in vec3 FragPos;
out vec4 fragColor;

float specularStrength = 0.5;

void main() {

    vec3 norm = normalize(Normal);

    vec3 lightDiff = lightPos - FragPos;
    float lightDistance = length(lightDiff);
    vec3 lightDir = normalize(lightDiff);

    float attenuation = 1.0 / (1 + 0.22 * lightDistance + 0.20 * (lightDistance * lightDistance));

    float diff = max(dot(norm, lightDir), 0.0);

    vec3 diffuseLight = diff * lightColor;

    vec3 viewDir = normalize(viewPos - FragPos);
    vec3 halfwayDir = normalize(lightDir + viewDir);

    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(norm, halfwayDir), 0.0), 8);
    vec3 specular = specularStrength * spec * lightColor;

    fragColor = texture(diffuse, TexCoord) * vec4(diffuseLight + specular, 1.0f) * attenuation;
}
