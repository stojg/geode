
uniform float x_varianceMin;
uniform float x_lightBleedReductionAmount;
in vec4 FragPosLightSpace;
uniform sampler2D x_shadowMap;

float linstep(float low, float high, float v)
{
	return clamp((v-low)/(high-low), 0.0, 1.0);
}

float sampleVarianceShadowMap(sampler2D shadowMap, vec2 coords, float compare, float varianceMin, float lightBleedReductionAmount)
{
    // return step(compare, texture(shadowMap, coords.xy).r);
    vec2 moments = texture(shadowMap, coords).xy;
	float p = step(compare, moments.x);

	float variance = max(moments.y - moments.x * moments.x, varianceMin);

	float d = compare - moments.x;
	float pMax = linstep(lightBleedReductionAmount, 1.0, variance / (variance + d*d));

	return min(max(p, pMax), 1.0);
}

float ShadowCalculation(vec4 fragPosLightSpace, vec3 normal, vec3 lightDir, float varianceMin, float lightBleedReductionAmount)
{
    // perform perspective divide, since it's not done automatically for us
    vec3 projCoords = fragPosLightSpace.xyz / fragPosLightSpace.w;
    // transform from [-0.5,0.5] to [0,1] range so we can use it for sampling
    projCoords = projCoords * 0.5 + 0.5;

    // dont shadow things outside the light frustrum far plane
    if(projCoords.z > 1.0) {
        return 1.0;
    }

    return sampleVarianceShadowMap(x_shadowMap, projCoords.xy, projCoords.z, varianceMin, lightBleedReductionAmount);
}
