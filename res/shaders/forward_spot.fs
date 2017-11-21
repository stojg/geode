#version 410 core

uniform sampler2D diffuse;
uniform mat4 view;

in vec2 TexCoord;
in vec3 LightPos;
in vec3 Normal;
in vec3 FragPos;

struct Attenuation
{
    float constant;
    float linear;
    float exponent;
};

struct PointLight
{
    vec3 color;
    vec3 position;
    Attenuation atten;
};

struct SpotLight
{
    PointLight pointLight;
    vec3 direction;
    float cutoff;
};

uniform SpotLight spotLight;

out vec4 fragColor;

float specularStrength = 0.5;

void main() {

    vec3 norm = normalize(Normal);

    vec3 lightDiff = LightPos - FragPos;

    float lightDistance = length(lightDiff);
    vec3 lightDir = normalize(lightDiff);

    vec3 spotlightDirection = (view * vec4(spotLight.direction,0)).xyz;
    float theta = dot(lightDir, normalize(-spotlightDirection));
    if(theta < spotLight.cutoff) {
        fragColor = vec4(0,0,0,0);
        return;
    }

    float attenuation = 1.0 / (spotLight.pointLight.atten.constant + spotLight.pointLight.atten.linear * lightDistance + spotLight.pointLight.atten.exponent * (lightDistance * lightDistance));

    float diff = max(dot(norm, lightDir), 0.0);

    vec3 diffuseLight = diff * spotLight.pointLight.color;

    vec3 viewDir = normalize(-FragPos);
    vec3 halfwayDir = normalize(lightDir + viewDir);

    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(norm, halfwayDir), 0.0), 128);
    vec3 specular = specularStrength * spec * spotLight.pointLight.color;

    fragColor = texture(diffuse, TexCoord) * vec4(diffuseLight + specular, 1.0f) * attenuation;
}
