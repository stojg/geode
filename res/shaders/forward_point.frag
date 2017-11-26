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


in vec2 TexCoord;
in vec3 LightPos;
in vec3 Normal;
in vec3 ModelViewPos;

out vec4 fragColor;

const float specularStrength = 0.5;

uniform sampler2D diffuse;
uniform PointLight pointLight;


void main() {

    vec3 norm = Normal;
    vec3 color = pointLight.base.color;

    vec3 lightDiff = LightPos - ModelViewPos;
    float lightDistance = length(lightDiff);

    vec3 lightDir = normalize(lightDiff);

    float attenuation = 1.0 / (pointLight.atten.constant + pointLight.atten.linear * lightDistance + pointLight.atten.exponent * (lightDistance * lightDistance));

    float diff = max(dot(norm, lightDir), 0.0);

    vec3 diffuseLight = diff * color;

    vec3 halfwayDir = normalize(lightDir - normalize(ModelViewPos));
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(norm, halfwayDir), 0.0), 128);
    vec3 specular = specularStrength * spec * color;

    float shadow = 1.0;

    fragColor = texture(diffuse, TexCoord);
    fragColor *= vec4((diffuseLight + specular), 1.0f);
    fragColor *= attenuation;
    fragColor *= shadow;
}
