#version 410 core

#include "light.frag"
#include "light_shadow.frag"

uniform DirectionalLight directionalLight;

void main() {

    vec3 lightDir = LightPos;

    vec3 color = directionalLight.base.color;

    vec3 diffuseLight = diffuseCalc(Normal, lightDir, color);
    vec3 specular = specularCalc(Normal, lightDir, color, specularStrength);

    // calculate shadow
    float shadow = 1.0;

    if (x_varianceMin != 0.0) {
        shadow = ShadowCalculation(FragPosLightSpace, Normal, lightDir, x_varianceMin, x_lightBleedReductionAmount);
    }

    fragColor = texture(diffuse, TexCoord);
    fragColor *= vec4((diffuseLight + specular), 1.0f);
    fragColor *= shadow;
}
