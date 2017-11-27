#version 410 core

#include "light.frag"
#include "light_shadow.frag"

uniform SpotLight spotLight;
uniform mat4 view;

bool isOutsideSpotLight(mat4 view, vec3 lightDir, vec3 direction, float cutOff) {

    vec3 viewDirection = (view * vec4(direction, 0)).xyz;

    float theta = dot(lightDir, normalize(-viewDirection));
    return theta < spotLight.cutoff;
}

void main() {

    vec3 lightDiff = LightPos - ModelViewPos;
    float lightDistance = length(lightDiff);
    vec3 ligthDirection = normalize(lightDiff);

    if(isOutsideSpotLight(view, ligthDirection, spotLight.direction, spotLight.cutoff)) {
        fragColor = vec4(0);
        return;
    }

    vec3 color = spotLight.pointLight.base.color;
    float attenuation = attenuationCalc(lightDistance, spotLight.pointLight.atten);
    vec3 diffuseLight = diffuseCalc(Normal, ligthDirection, color);
    vec3 specular = specularCalc(Normal, ligthDirection, color, specularStrength);

    // calculate shadow
    float shadow = 1.0;
    if (x_varianceMin != 0.0) {
        shadow = ShadowCalculation(FragPosLightSpace, Normal, ligthDirection, x_varianceMin, x_lightBleedReductionAmount);
    }

    fragColor = texture(diffuse, TexCoord);
    fragColor *= vec4((diffuseLight + specular), 1.0f);
    fragColor *= attenuation;
    fragColor *= shadow;
}
