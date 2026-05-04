//+vert
#version 330 core

layout (location = 0) in vec3 a_Pos;
layout (location = 4) in vec2 a_UV0;
// layout (location = 0) in vec3 a_Pos;
// layout (location = 1) in vec3 a_Norm;
// layout (location = 2) in vec2 a_UV0;

out vec2 v_UV0;

uniform mat4 u_PMatrix;
uniform mat4 u_VMatrix;
uniform mat4 u_MMatrix;

void main() {
    gl_Position = u_PMatrix * u_VMatrix * u_MMatrix * vec4(a_Pos.xyz, 1.0);
    v_UV0 = a_UV0;
}
//+frag
#version 330 core
precision mediump float;

out vec4 f_Color;

in vec2 v_UV0;

uniform sampler2D u_DiffTexture;
/*
uniform vec3 u_DirLightDir[2];
uniform vec3 u_DirLightColor[2];
uniform float u_DirLightPower[2];

uniform vec3 u_PtLightPos[8];
uniform vec3 u_PtLightColor[8];
uniform float u_PtLightPower[8];
*/
void main() {
    f_Color = texture2D(u_DiffTexture, v_UV0);
}