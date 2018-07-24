const vec3 fogColor = vec3(0.5,0.6,0.7);
const float FogDensity = 0.0000001;

vec3 fogCalc(vec3 Lo, vec3 vp) {
    float dist = 0;
    dist = length(vp);

    float fogAmount = 1.0 - exp(-dist *dist *dist * FogDensity);
    fogAmount = clamp( fogAmount, 0.0, 1.0 );
    return mix(Lo, fogColor, fogAmount);
}
