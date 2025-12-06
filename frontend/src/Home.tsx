import { APITester } from "./APITester";
import "./index.css";

import logo from "./logo.svg";
import reactLogo from "./react.svg";

export function Home() {
  return (
    <div className="max-w-7xl mx-auto p-8 text-center relative z-10">
      <div className="flex justify-center items-center gap-8 mb-8">
        <img
          src={logo}
          alt="Bun Logo"
          className="h-24 p-6 transition-all duration-300 hover:drop-shadow-[0_0_2em_#646cffaa] scale-120"
        />
        <img
          src={reactLogo}
          alt="React Logo"
          className="h-24 p-6 transition-all duration-300 hover:drop-shadow-[0_0_2em_#61dafbaa] animate-[spin_20s_linear_infinite]"
        />
      </div>

      <h1 className="text-5xl font-bold my-4 leading-tight">Bun + React</h1>
      <p>
        Edit <code className="bg-[#1a1a1a] px-2 py-1 rounded font-mono">src/App.tsx</code> and save to test HMR
      </p>

      <div className="my-8 flex gap-4 justify-center flex-wrap">
        <a
          href="/chat"
          className="inline-block bg-[#646cff] hover:bg-[#535bf2] text-white font-medium py-2 px-6 rounded-lg transition-colors duration-200"
        >
          Start Chat
        </a>
        <a
          href="/preferences"
          className="inline-block bg-[#2a2a2a] hover:bg-[#333] text-white font-medium py-2 px-6 rounded-lg border border-[#444] transition-colors duration-200"
        >
          Preferences
        </a>
        <a
          href="/login"
          className="inline-block bg-[#2a2a2a] hover:bg-[#333] text-white font-medium py-2 px-6 rounded-lg border border-[#444] transition-colors duration-200"
        >
          Login
        </a>
      </div>

      <APITester />
    </div>
  );
}

export default Home;
