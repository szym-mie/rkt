//+vert
#version 330 core

layout (location = 0) in vec2 a_Pos;

out vec4 v_Pos;

void main() {
    gl_Position = vec4(a_Pos.xy, 1.0, 1.0);
    v_Pos = vec4(a_Pos.xy, 0.0, 1.0);
}
//+frag
#version 330 core
precision mediump float;

out vec4 f_Color;

in vec4 v_Pos;

uniform mat4 u_PMatrix;
uniform mat4 u_VMatrix;

uniform samplerCube u_DiffTexture;

void main() {
    mat4 vd = u_VMatrix;
    vd[3] = vec4(0.0, 0.0, 0.0, 1.0);
    vec4 back = inverse(u_PMatrix * vd) * v_Pos;
    vec3 uv = normalize(back.xyz / back.w);
    float toSun = distance(uv, normalize(vec3(0.7,0.3,0.5)));
    float inSun = max(toSun * 100.0 - 0.5, 0.0);
    float z = max(uv.z, 0.0);
    vec3 blue = mix(vec3(0.0, 0.8, 1.0), vec3(0.0, 0.4, 1.0), sqrt(z));
    vec3 sun = mix(vec3(1.00, 1.00, 1.00), vec3(0.98, 1.00, 0.65), inSun);
    vec3 color = mix(sun, blue, min(max(inSun * inSun, 0.0), 1.0));
    f_Color = vec4(color, 1.0);
}