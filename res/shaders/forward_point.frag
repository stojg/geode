#version 410 core

#include "light.frag"

uniform PointLight pointLight;

void main() {

    vec3 lightDiff = LightPos - ModelViewPos;
    float lightDistance = length(lightDiff);
    vec3 lightDir = normalize(lightDiff);

    vec3 color = pointLight.base.color;

    float attenuation = attenuationCalc(lightDistance, pointLight.atten);
    vec3 diffuseLight = diffuseCalc(Normal, lightDir, color);
    vec3 specular = specularCalc(Normal, lightDir, color, specularStrength);

    fragColor = texture(diffuse, TexCoord);
    fragColor *= vec4((diffuseLight + specular), 1.0f);
    fragColor *= attenuation;
}
