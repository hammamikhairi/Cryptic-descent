#version 330

// Input vertex attributes (from vertex shader)
in vec2 fragTexCoord;
in vec4 fragColor;

// Output fragment color
out vec4 finalColor;

// Uniform inputs
uniform sampler2D texture0;
uniform vec2 resolution;
uniform vec2 lightPos[16];
uniform float lightRadius[16];
uniform int lightCount;
uniform float decayFactor;
uniform float time;
uniform int lightModes[16];

void main()
{
    // Early exit if no lights
    if (lightCount <= 0) {
        finalColor = vec4(0.0);
        return;
    }

    vec2 pixelPos = gl_FragCoord.xy;
    vec4 color = vec4(0.0);
    float totalIntensity = 0.0;
    
    // Process only the closest lights
    for(int i = 0; i < lightCount; i++) {
        vec2 delta = pixelPos - lightPos[i];
        float distance = length(delta);
        
        // Early distance check
        if(distance >= lightRadius[i]) continue;
        
        float intensity = exp(-decayFactor * (distance/lightRadius[i]));
        
        // Apply effects based on mode (simplified)
        switch(lightModes[i]) {
            case 1: // Shimmer
                intensity *= 0.8 + 0.2 * sin(time * 2.0);
                break;
            case 2: // Pulse
                intensity *= 0.8 + 0.2 * sin(time * 6.0);
                break;
            // Add more cases as needed, but keep them simple
        }
        
        totalIntensity += min(intensity * 1.2, 1.0);
    }
    
    // Clamp final intensity
    finalColor = vec4(vec3(min(totalIntensity, 1.0)), 1.0);
}