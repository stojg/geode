#version 410 core

struct Attenuation
{
    float constant;
    float linear;
    float exponent;
};

struct BaseLight
{
    vec3 color;
};

struct DirectionalLight
{
    BaseLight base;
    vec3 direction;
};

struct PointLight
{
    BaseLight base;
    Attenuation atten;
    vec3 position;
};

struct SpotLight
{
    PointLight pointLight;
    vec3 direction;
    float cutoff;
};

uniform sampler2D diffuse;

in vec2 TexCoord;
in vec3 LightPos;
in vec3 Normal;
in vec3 FragPos;

uniform PointLight pointLight;

out vec4 fragColor;

float specularStrength = 0.5;

void main() {

    vec3 norm = normalize(Normal);

    vec3 lightDiff = LightPos - FragPos;
    float lightDistance = length(lightDiff);
    vec3 lightDir = normalize(lightDiff);

    float attenuation = 1.0 / (pointLight.atten.constant + pointLight.atten.linear * lightDistance + pointLight.atten.exponent * (lightDistance * lightDistance));

    float diff = max(dot(norm, lightDir), 0.0);

    vec3 diffuseLight = diff * pointLight.base.color;

    vec3 viewDir = normalize(-FragPos);
    vec3 halfwayDir = normalize(lightDir + viewDir);
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(norm, halfwayDir), 0.0), 128);
    vec3 specular = specularStrength * spec * pointLight.base.color;

    fragColor = texture(diffuse, TexCoord) * vec4(diffuseLight + specular, 1.0f) * attenuation;
}
