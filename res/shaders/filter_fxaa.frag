#version 410 core

#define FXAA_SEARCH_ACCELERATION 0

out vec4 FragColor;

in vec2 TexCoords;

uniform sampler2D x_filterTexture;

uniform vec2 u_texelStep = vec2(1/800, 1/600);
uniform int u_showEdges = 1;

uniform float u_lumaThreshold = 0.1f; // 0.45 - 0.8 - 0.6

// The minimum amount of local contrast required to apply algorithm.
// 1/3 – too little
// 1/4 – low quality
// 1/8 – high quality
// 1/16 – overkill
uniform float u_mulReduce = 1 / 8.0f; // 8.0f

// Trims the algorithm from processing darks.
// 1/32 – visible limit
// 1/16 – high quality
// 1/12 – upper limit (start of visible unfiltered edges)
uniform float u_minReduce = 1 / 12.0f; // 1 / 128.0f


uniform float u_maxSpan = 8.0f; // 8.0f;

uniform uint rt_w = 800;
uniform uint rt_h = 600;

// see FXAA
// http://developer.download.nvidia.com/assets/gamedev/files/sdk/11/FXAA_WhitePaper.pdf
// http://iryoku.com/aacourse/downloads/09-FXAA-3.11-in-15-Slides.pdf
// http://horde3d.org/wiki/index.php5?title=Shading_Technique_-_FXAA
// https://github.com/Asmodean-/PsxFX/blob/master/PsxFX/gpuPeteOGL2.slf

// As an optimization, luminance is estimated strictly from Red and Green channels
// using a single fused multiply add operation. In practice pure blue aliasing rarely
// appears in typical game content.
float FxaaLuma(vec3 val) {
 return val.y * (0.587/0.299) + val.x;
}

#define FxaaInt2 ivec2
#define FxaaFloat2 vec2
#define FxaaTexLod0(t, p) textureLod(t, p, 0.0)
#define FxaaTexOff(t, p, o, r) textureLodOffset(t, p, 0.0, o)

#define FXAA_REDUCE_MIN   (1.0/128.0)
#define FXAA_REDUCE_MUL   (1.0/8.0)
#define FXAA_SPAN_MAX     8.0

// Output of FxaaVertexShader interpolated across screen.
// Input texture.
// Constant {1.0/frameWidth, 1.0/frameHeight}.
vec3 FxaaPixelShader(vec2 posPos, sampler2D tex, vec2 rcpFrame)
{

    vec3 rgbNW = FxaaTexLod0(tex, posPos).xyz;
    vec3 rgbNE = FxaaTexOff(tex, posPos, FxaaInt2(1,0), rcpFrame.xy).xyz;
    vec3 rgbSW = FxaaTexOff(tex, posPos, FxaaInt2(0,1), rcpFrame.xy).xyz;
    vec3 rgbSE = FxaaTexOff(tex, posPos, FxaaInt2(1,1), rcpFrame.xy).xyz;
    vec3 rgbM  = FxaaTexLod0(tex, posPos.xy).xyz;

    vec3 luma = vec3(0.299, 0.587, 0.114);
    float lumaNW = dot(rgbNW, luma);
    float lumaNE = dot(rgbNE, luma);
    float lumaSW = dot(rgbSW, luma);
    float lumaSE = dot(rgbSE, luma);
    float lumaM  = dot(rgbM,  luma);

    float lumaMin = min(lumaM, min(min(lumaNW, lumaNE), min(lumaSW, lumaSE)));
    float lumaMax = max(lumaM, max(max(lumaNW, lumaNE), max(lumaSW, lumaSE)));

    vec2 dir;
    dir.x = -((lumaNW + lumaNE) - (lumaSW + lumaSE));
    dir.y =  ((lumaNW + lumaSW) - (lumaNE + lumaSE));

    float dirReduce = max((lumaNW + lumaNE + lumaSW + lumaSE) * (0.25 * FXAA_REDUCE_MUL), FXAA_REDUCE_MIN);
    float rcpDirMin = 1.0/(min(abs(dir.x), abs(dir.y)) + dirReduce);
    dir = min(FxaaFloat2( FXAA_SPAN_MAX,  FXAA_SPAN_MAX), max(FxaaFloat2(-FXAA_SPAN_MAX, -FXAA_SPAN_MAX), dir * rcpDirMin)) * rcpFrame.xy;

    vec3 a = FxaaTexLod0(tex, posPos.xy + dir * (1.0/3.0 - 0.5)).xyz + FxaaTexLod0(tex, posPos.xy + dir * (2.0/3.0 - 0.5)).xyz;
    vec3 rgbA = (1.0/2.0) * a;
    vec3 b = (FxaaTexLod0(tex, posPos.xy + dir * (0.0/3.0 - 0.5)).xyz + FxaaTexLod0(tex, posPos.xy + dir * (3.0/3.0 - 0.5)).xyz);
    vec3 rgbB = rgbA * (1.0/2.0) + (1.0/4.0) * b;

    float lumaB = dot(rgbB, luma);

    if((lumaB < lumaMin) || (lumaB > lumaMax)) {
        return rgbA;
    }
    return rgbB;
}


void main(void)
{
    vec2 rcpFrame = vec2(1.0/rt_w, 1.0/rt_h);
    FragColor.rgb = FxaaPixelShader(TexCoords, x_filterTexture, rcpFrame);
    FragColor.w = 1;
}
