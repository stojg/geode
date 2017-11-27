
#include "light.glh"

in vec2 TexCoord;
in vec3 LightPos;
in vec3 Normal;
in vec3 ModelViewPos;

out vec4 fragColor;

const float specularStrength = 0.5;

uniform sampler2D diffuse;

vec3 diffuseCalc(vec3 norm, vec3 lightDirection, vec3 color) {
    return max(dot(norm, lightDirection), 0.0) * color;
}

vec3 specularCalc(vec3 norm, vec3 lightDirection, vec3 color, float strength) {
    vec3 halfwayDir = normalize(lightDirection - normalize(ModelViewPos));
    vec3 reflectDir = reflect(-lightDirection, norm);
    float spec = pow(max(dot(norm, halfwayDir), 0.0), 128);
    return strength * spec * color;
}

float attenuationCalc(float lightDistance, Attenuation atten) {
    return 1.0 / (atten.constant + atten.linear * lightDistance + atten.exponent * (lightDistance * lightDistance));
}
