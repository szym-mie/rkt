//+vert
#version 330 core

layout (location = 0) in vec3 a_Pos;
layout (location = 1) in vec3 a_Norm;
layout (location = 4) in vec2 a_UV0;

out vec3 v_FragPos;
out vec3 v_Norm;
out vec2 v_UV0;

uniform mat4 u_PMatrix;
uniform mat4 u_VMatrix;
uniform mat4 u_MMatrix;

void main() {
    gl_Position = u_PMatrix * u_VMatrix * u_MMatrix * vec4(a_Pos, 1.0);
    v_Norm = mat3(u_MMatrix) * a_Norm;
    v_FragPos = vec3(u_MMatrix * vec4(a_Pos, 1.0));
    v_UV0 = a_UV0;
}
//+frag
#version 330 core
precision mediump float;

out vec4 f_Color;

in vec3 v_FragPos;
in vec3 v_Norm;
in vec2 v_UV0;

uniform sampler2D u_DiffTexture;
uniform vec3 u_AmbLightColor;
uniform vec3 u_DirLightDir[2];
uniform vec3 u_DirLightColor[2];
/*
uniform vec3 u_PtLightPos[8];
uniform vec3 u_PtLightColor[8];
uniform vec3 u_PtLightPower[8];
*/

vec3 calcDirLight(vec3 dir, vec3 color, vec3 norm) {
    float diff = max(dot(dir, norm), 0.0);
    vec3 diffuse = color * diff * vec3(texture(u_DiffTexture, v_UV0));
    return diffuse;
}

vec3 calcPtLight(vec3 pos, vec3 color, float linFall, float quaFall, vec3 norm) {
    vec3 dir = normalize(pos - v_FragPos);
    float diff = max(dot(dir, norm), 0.0);
    float dist = length(pos - v_FragPos);
    float atten = 1.0 / (1.0 + linFall * dist + quaFall * (dist * dist));
    vec3 diffuse = color * diff * vec3(texture(u_DiffTexture, v_UV0));
    diffuse *= atten;
    return diffuse;
}

void main() {
    vec3 norm = normalize(v_Norm);
    vec3 color = u_AmbLightColor * vec3(texture(u_DiffTexture, v_UV0));
    for (int i = 0; i < 2; i++) {
        color += calcDirLight(u_DirLightDir[i], u_DirLightColor[i], norm);
    }
    f_Color = vec4(color, 1.0);
}