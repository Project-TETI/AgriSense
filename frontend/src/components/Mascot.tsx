'use client';

import React from 'react';

interface MascotProps {
  isSmiling?: boolean;
}

export const Mascot: React.FC<MascotProps> = ({ isSmiling = false }) => {
  return (
    <div className="relative w-48 h-48 sm:w-56 sm:h-56 flex items-center justify-center select-none">
      {/* Encouraging Speech Bubble Pill */}
      <div 
        className={`absolute -top-2 right-0 bg-white/95 px-3 py-1 rounded-full shadow-md border border-white text-[11px] font-bold text-forest-800 z-20 flex items-center gap-1.5 transition-all duration-300 ${
          isSmiling ? 'scale-105 shadow-lg border-emerald-300/60' : ''
        }`}
      >
        <span>Siap bantu!</span> 🌿
      </div>

      {/* Mascot SVG with floating animation */}
      <svg 
        className="w-full h-full filter drop-shadow-xl" 
        fill="none" 
        viewBox="0 0 280 280" 
        xmlns="http://www.w3.org/2000/svg"
      >
        <defs>
          {/* 3D Organic Leaf Body Base Gradients */}
          <radialGradient id="charBodyShade" cx="35%" cy="30%" r="70%">
            <stop offset="0%" stopColor="#A8F572" />
            <stop offset="28%" stopColor="#5CD368" />
            <stop offset="65%" stopColor="#27A454" />
            <stop offset="92%" stopColor="#136E39" />
            <stop offset="100%" stopColor="#0B4725" />
          </radialGradient>
          <radialGradient id="charHighlightSpecular" cx="32%" cy="22%" r="40%">
            <stop offset="0%" stopColor="#FFFFFF" stopOpacity="0.8" />
            <stop offset="40%" stopColor="#E2FFC6" stopOpacity="0.4" />
            <stop offset="100%" stopColor="#5CD368" stopOpacity="0" />
          </radialGradient>
          <linearGradient id="charRimLight" x1="85%" y1="10%" x2="15%" y2="90%">
            <stop offset="0%" stopColor="#BAFF8A" stopOpacity="0.65" />
            <stop offset="50%" stopColor="#38C472" stopOpacity="0" />
            <stop offset="100%" stopColor="#09381E" stopOpacity="0.7" />
          </linearGradient>
          {/* Ambient Shadow Underneath Body */}
          <radialGradient id="dropFloorShadow" cx="50%" cy="50%" r="50%">
            <stop offset="0%" stopColor="#0D3B24" stopOpacity="0.35" />
            <stop offset="60%" stopColor="#0D3B24" stopOpacity="0.15" />
            <stop offset="100%" stopColor="#0D3B24" stopOpacity="0" />
          </radialGradient>
          {/* Sprout Leaves Shading */}
          <linearGradient id="sproutMainGrad" x1="20%" y1="10%" x2="80%" y2="90%">
            <stop offset="0%" stopColor="#D3FFAA" />
            <stop offset="30%" stopColor="#6DE867" />
            <stop offset="100%" stopColor="#1E833B" />
          </linearGradient>
          <linearGradient id="sproutStemGrad" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor="#5CD368" />
            <stop offset="100%" stopColor="#155E2B" />
          </linearGradient>
          {/* Researcher Glasses Metallic Frames */}
          <linearGradient id="glassesFrameGrad" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor="#54E3C0" />
            <stop offset="45%" stopColor="#1E6B60" />
            <stop offset="80%" stopColor="#0D3D37" />
            <stop offset="100%" stopColor="#35B398" />
          </linearGradient>
          {/* Lens Sheen Glass Gradient */}
          <linearGradient id="lensGlassGrad" x1="20%" y1="15%" x2="85%" y2="85%">
            <stop offset="0%" stopColor="#FFFFFF" stopOpacity="0.75" />
            <stop offset="25%" stopColor="#E3FAFF" stopOpacity="0.45" />
            <stop offset="65%" stopColor="#8CE7FF" stopOpacity="0.1" />
            <stop offset="100%" stopColor="#1E6B60" stopOpacity="0.15" />
          </linearGradient>
          {/* Magnifying Glass Elements */}
          <linearGradient id="magRimGrad" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor="#E8FDF5" />
            <stop offset="40%" stopColor="#2B7A70" />
            <stop offset="100%" stopColor="#103833" />
          </linearGradient>
          <radialGradient id="magLensReflect" cx="35%" cy="30%" r="60%">
            <stop offset="0%" stopColor="#FFFFFF" stopOpacity="0.85" />
            <stop offset="50%" stopColor="#ABF5E5" stopOpacity="0.45" />
            <stop offset="100%" stopColor="#184F48" stopOpacity="0.25" />
          </radialGradient>
          <linearGradient id="magHandleGrad" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor="#D1A153" />
            <stop offset="40%" stopColor="#8E6221" />
            <stop offset="100%" stopColor="#543A11" />
          </linearGradient>
          {/* Soft Kawaii Blush */}
          <radialGradient id="kawaiiBlush" cx="50%" cy="50%" r="50%">
            <stop offset="0%" stopColor="#FF6B8B" stopOpacity="0.75" />
            <stop offset="60%" stopColor="#FF8FA6" stopOpacity="0.35" />
            <stop offset="100%" stopColor="#FF8FA6" stopOpacity="0" />
          </radialGradient>
          {/* 3D Paws/Arms Gradient */}
          <radialGradient id="pawGrad" cx="35%" cy="30%" r="65%">
            <stop offset="0%" stopColor="#9BF477" />
            <stop offset="50%" stopColor="#42C466" />
            <stop offset="100%" stopColor="#187037" />
          </radialGradient>
        </defs>

        {/* Ground Contact Shadow with breathing animation */}
        <ellipse 
          cx="138" 
          cy="250" 
          rx="74" 
          ry="14" 
          fill="url(#dropFloorShadow)" 
          className="animate-mascot-shadow"
        />

        {/* Floating Mascot Group */}
        <g className="animate-mascot-float">
          {/* Sprout Stem & Curly Top Leaf with subtle sway */}
          <g id="sprout-crest" className="animate-sprout-sway">
            {/* Curving Stem */}
            <path d="M138 72 C138 52 144 32 152 24 C153 23 155 24 154 26 C147 38 144 54 142 72 Z" fill="url(#sproutStemGrad)" />
            {/* Main Sprout Leaf (Right) */}
            <path d="M152 25 C176 16 194 34 178 52 C164 64 148 50 152 25 Z" fill="url(#sproutMainGrad)" filter="drop-shadow(0px 2px 3px rgba(10,50,25,0.25))" />
            <path d="M154 28 Q166 42 174 48" stroke="#FFFFFF" strokeWidth="1.5" strokeLinecap="round" strokeOpacity="0.6" fill="none" />
            {/* Second Baby Sprout Leaf (Left) */}
            <path d="M140 38 C120 28 112 44 124 56 C134 64 144 54 140 38 Z" fill="url(#sproutMainGrad)" opacity="0.95" />
            <path d="M138 40 Q130 48 126 53" stroke="#FFFFFF" strokeWidth="1.2" strokeLinecap="round" strokeOpacity="0.55" fill="none" />
          </g>

          {/* MAIN 3D ASYMMETRIC TEARDROP BOTANICAL BODY */}
          <g id="mascot-body-group">
            {/* Teardrop Leaf Contour */}
            <path d="M138 68 C168 69 196 92 208 124 C222 160 220 196 201 223 C184 246 156 252 135 252 C108 252 76 242 62 216 C46 186 48 146 68 114 C84 86 112 67 138 68 Z" fill="url(#charBodyShade)" />
            {/* Rim Light Layer for 3D Volume Edge */}
            <path d="M138 68 C168 69 196 92 208 124 C222 160 220 196 201 223 C184 246 156 252 135 252 C108 252 76 242 62 216 C46 186 48 146 68 114 C84 86 112 67 138 68 Z" fill="url(#charRimLight)" />
            {/* Big Specular Volume Highlight */}
            <ellipse cx="116" cy="112" rx="54" ry="40" transform="rotate(-18 116 112)" fill="url(#charHighlightSpecular)" />
            {/* Gloss Reflection Dot on Top Shoulder */}
            <ellipse cx="102" cy="90" rx="14" ry="6" transform="rotate(-22 102 90)" fill="#FFFFFF" fillOpacity="0.6" />

            {/* Leaf Vein Ribs */}
            <path d="M138 74 C140 108 138 178 134 244" stroke="#78EA82" strokeWidth="3" strokeLinecap="round" strokeOpacity="0.45" fill="none" />
            <path d="M138 120 C158 130 182 144 196 154" stroke="#78EA82" strokeWidth="2" strokeLinecap="round" strokeOpacity="0.35" fill="none" />
            <path d="M136 142 C112 152 92 162 76 170" stroke="#78EA82" strokeWidth="2" strokeLinecap="round" strokeOpacity="0.35" fill="none" />
            <path d="M136 176 C154 186 176 196 188 204" stroke="#78EA82" strokeWidth="2" strokeLinecap="round" strokeOpacity="0.35" fill="none" />
          </g>

          {/* Kawaii Rosy Cheeks */}
          <g 
            id="blush-cheeks"
            className="transition-transform duration-300 ease-in-out"
            style={{
              transform: isSmiling ? 'scale(1.15)' : 'scale(1)',
              transformOrigin: '136px 174px',
            }}
          >
            <ellipse cx="84" cy="174" rx="14" ry="9" transform="rotate(-5 84 174)" fill="url(#kawaiiBlush)" />
            <ellipse cx="188" cy="174" rx="14" ry="9" transform="rotate(5 188 174)" fill="url(#kawaiiBlush)" />
            {/* Freckle stars */}
            <circle cx="81" cy="172" r="1.2" fill="#FFFFFF" opacity="0.7" />
            <circle cx="88" cy="175" r="1" fill="#FFFFFF" opacity="0.7" />
            <circle cx="184" cy="174" r="1" fill="#FFFFFF" opacity="0.7" />
            <circle cx="191" cy="172" r="1.2" fill="#FFFFFF" opacity="0.7" />
          </g>

          {/* Expressive Kawaii Eyes */}
          <g id="mascot-eyes">
            {/* Left Eye */}
            <g id="left-eye">
              <ellipse cx="104" cy="150" rx="12.5" ry="15" fill="#0A2619" />
              <ellipse cx="104" cy="152" rx="10" ry="11" fill="#175438" />
              <circle cx="99" cy="143" r="5.2" fill="#FFFFFF" />
              <circle cx="108" cy="155" r="2.6" fill="#FFFFFF" />
              <circle cx="102" cy="159" r="1.2" fill="#D2FFE0" />
              {/* Extra twinkle sparkle on smile */}
              <circle 
                cx="97" 
                cy="151" 
                r="1.8" 
                fill="#FFFFFF" 
                className="transition-opacity duration-300"
                style={{ opacity: isSmiling ? 0.9 : 0 }} 
              />
            </g>
            {/* Right Eye */}
            <g id="right-eye">
              <ellipse cx="168" cy="150" rx="12.5" ry="15" fill="#0A2619" />
              <ellipse cx="168" cy="152" rx="10" ry="11" fill="#175438" />
              <circle cx="163" cy="143" r="5.2" fill="#FFFFFF" />
              <circle cx="172" cy="155" r="2.6" fill="#FFFFFF" />
              <circle cx="166" cy="159" r="1.2" fill="#D2FFE0" />
              {/* Extra twinkle sparkle on smile */}
              <circle 
                cx="161" 
                cy="151" 
                r="1.8" 
                fill="#FFFFFF" 
                className="transition-opacity duration-300"
                style={{ opacity: isSmiling ? 0.9 : 0 }} 
              />
            </g>
          </g>

          {/* Mouth and Smile Expressions with smooth 200-300ms transition */}
          <g id="mouth-and-smile">
            {/* Default gentle sweet smile */}
            <g 
              className="transition-all duration-300 ease-in-out"
              style={{
                opacity: isSmiling ? 0 : 1,
                transform: isSmiling ? 'scale(0.85)' : 'scale(1)',
                transformOrigin: '136px 172px',
              }}
            >
              <path d="M127 168 C131 178 141 178 145 168" stroke="#0C291B" strokeWidth="3.5" strokeLinecap="round" fill="#991E3B" />
              <path d="M131 172 C134 176 138 176 141 172" fill="#FF7A95" />
            </g>

            {/* Hovered responsive beaming happy smile */}
            <g 
              className="transition-all duration-300 ease-in-out"
              style={{
                opacity: isSmiling ? 1 : 0,
                transform: isSmiling ? 'scale(1)' : 'scale(0.85)',
                transformOrigin: '136px 172px',
              }}
            >
              {/* Joyful open mouth curve */}
              <path 
                d="M124 167 C125 183 147 183 148 167 Z" 
                stroke="#0C291B" 
                strokeWidth="3.2" 
                strokeLinecap="round" 
                fill="#8C132E" 
              />
              {/* Sweet tongue in smile */}
              <path 
                d="M128 175 C132 181 140 181 144 175 Z" 
                fill="#FF7A95" 
              />
              {/* Cheerful mouth lip line outline */}
              <path 
                d="M124 167 C125 183 147 183 148 167" 
                stroke="#0C291B" 
                strokeWidth="3.2" 
                strokeLinecap="round" 
                fill="none" 
              />
            </g>
          </g>

          {/* Researcher Round Glasses */}
          <g id="researcher-glasses" filter="drop-shadow(0px 3px 6px rgba(12,45,28,0.22))">
            {/* Bridge */}
            <path d="M124 148 Q136 142 148 148" stroke="url(#glassesFrameGrad)" strokeWidth="4" strokeLinecap="round" fill="none" />
            <path d="M124 148 Q136 142 148 148" stroke="#A5FDE8" strokeWidth="1.2" strokeLinecap="round" fill="none" opacity="0.8" />

            {/* Left Rim */}
            <circle cx="104" cy="150" r="23" fill="url(#lensGlassGrad)" stroke="url(#glassesFrameGrad)" strokeWidth="3.8" />
            <path d="M90 137 C95 131 106 130 114 133" stroke="#FFFFFF" strokeWidth="2.4" strokeLinecap="round" fill="none" opacity="0.85" />
            <circle cx="118" cy="140" r="1.8" fill="#FFFFFF" opacity="0.7" />

            {/* Right Rim */}
            <circle cx="168" cy="150" r="23" fill="url(#lensGlassGrad)" stroke="url(#glassesFrameGrad)" strokeWidth="3.8" />
            <path d="M154 137 C159 131 170 130 178 133" stroke="#FFFFFF" strokeWidth="2.4" strokeLinecap="round" fill="none" opacity="0.85" />
            <circle cx="182" cy="140" r="1.8" fill="#FFFFFF" opacity="0.7" />

            {/* Temple arms */}
            <path d="M81 149 L72 146" stroke="url(#glassesFrameGrad)" strokeWidth="3.2" strokeLinecap="round" />
            <path d="M191 149 L200 146" stroke="url(#glassesFrameGrad)" strokeWidth="3.2" strokeLinecap="round" />
          </g>

          {/* Paws and Magnifier */}
          <g id="hands-and-tools">
            {/* Left Waving Paw */}
            <g id="left-paw" transform="translate(56, 178) rotate(-18)">
              <ellipse cx="12" cy="12" rx="14" ry="10" fill="url(#pawGrad)" filter="drop-shadow(0 2px 4px rgba(10,40,20,0.2))" />
              <ellipse cx="9" cy="9" rx="6" ry="3" fill="#CFFFF0" opacity="0.6" />
            </g>

            {/* Right Arm & Magnifier Tool */}
            <g id="right-arm-and-magnifier">
              <g transform="translate(182, 142) rotate(-8)" filter="drop-shadow(2px 5px 8px rgba(12,45,28,0.28))">
                <rect x="38" y="54" width="10" height="36" rx="5" fill="url(#magHandleGrad)" stroke="#3D2706" strokeWidth="1.5" />
                <rect x="36" y="48" width="14" height="7" rx="2.5" fill="url(#glassesFrameGrad)" stroke="#184F48" strokeWidth="1" />
                <line x1="41" y1="58" x2="41" y2="84" stroke="#FFE4A0" strokeWidth="1.5" strokeLinecap="round" opacity="0.7" />

                <circle cx="28" cy="28" r="32" fill="url(#magRimGrad)" stroke="#0F3832" strokeWidth="2.5" />
                <circle cx="28" cy="28" r="26" fill="url(#magLensReflect)" />

                {/* Botanical micro-leaf inside lens */}
                <g transform="translate(18, 16)">
                  <path d="M8 18 C3 12 6 4 15 2 C22 9 20 18 12 21 Z" fill="#22C55E" stroke="#15803D" strokeWidth="1.5" />
                  <path d="M8 18 Q13 11 15 3" stroke="#DCFCE7" strokeWidth="1.2" fill="none" />
                  <circle cx="16" cy="7" r="2" fill="#FEF08A" />
                  <path d="M19 4 L20 1 L21 4 L24 5 L21 6 L20 9 L19 6 L16 5 Z" fill="#FEF08A" opacity="0.9" />
                </g>

                <path d="M12 16 C19 8 35 7 42 13" stroke="#FFFFFF" strokeWidth="3.8" strokeLinecap="round" fill="none" opacity="0.9" />
                <circle cx="14" cy="38" r="3" fill="#FFFFFF" opacity="0.7" />
              </g>

              {/* Hand paw gripping handle */}
              <circle cx="224" cy="198" r="12" fill="url(#pawGrad)" stroke="#135E2D" strokeWidth="1.5" filter="drop-shadow(0 2px 4px rgba(10,40,20,0.25))" />
              <ellipse cx="222" cy="194" rx="5" ry="3" fill="#CFFFF0" opacity="0.65" />
            </g>
          </g>

          {/* Floating AI Botanical Sparkles */}
          <g id="sparkles" opacity="0.85">
            <path d="M238 68 L240 60 L242 68 L250 70 L242 72 L240 80 L238 72 L230 70 Z" fill="#38A783" />
            <circle cx="240" cy="70" r="1.8" fill="#FFFFFF" />
            <path d="M48 94 L49.5 88 L51 94 L57 95.5 L51 97 L49.5 103 L48 97 L42 95.5 Z" fill="#5CD368" />
            <circle cx="49.5" cy="95.5" r="1.2" fill="#FFFFFF" />
          </g>
        </g>
      </svg>
    </div>
  );
};
