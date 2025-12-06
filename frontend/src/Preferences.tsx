import { useState, useEffect, FormEvent } from "react";
import { useNavigate } from "react-router-dom";

interface UserPreferences {
  dietary_restrictions: string;
  max_cooking_time: number;
}

export function Preferences() {
  const [dietaryRestrictions, setDietaryRestrictions] = useState("");
  const [maxCookingTime, setMaxCookingTime] = useState<number>(0);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const navigate = useNavigate();

  // Load preferences on mount
  useEffect(() => {
    loadPreferences();
  }, []);

  const loadPreferences = async () => {
    try {
      const response = await fetch("http://localhost:8080/api/preferences", {
        credentials: "include",
      });

      if (response.ok) {
        const data = await response.json();
        setDietaryRestrictions(data.dietary_restrictions || "");
        setMaxCookingTime(data.max_cooking_time || 0);
      } else if (response.status === 401) {
        // Not authenticated, redirect to login
        navigate("/login");
      } else {
        setError("Failed to load preferences");
      }
    } catch (err) {
      setError("Network error. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError("");
    setSuccess("");
    setSaving(true);

    try {
      const response = await fetch("http://localhost:8080/api/preferences", {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({
          dietary_restrictions: dietaryRestrictions,
          max_cooking_time: maxCookingTime,
        }),
      });

      const data = await response.json();

      if (response.ok) {
        setSuccess("Preferences saved successfully!");
        setTimeout(() => setSuccess(""), 3000);
      } else {
        setError(data.message || "Failed to save preferences");
      }
    } catch (err) {
      setError("Network error. Please try again.");
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-gray-400">Loading preferences...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen p-4">
      <div className="max-w-2xl mx-auto">
        {/* Header */}
        <div className="mb-6">
          <button
            onClick={() => navigate(-1)}
            className="text-[#646cff] hover:text-[#535bf2] mb-4 flex items-center gap-2 transition-colors"
          >
            ← Back
          </button>
          <h1 className="text-3xl font-bold mb-2">Preferences</h1>
          <p className="text-gray-400">
            Customize your meal planning experience
          </p>
        </div>

        {/* Form */}
        <div className="bg-[#1a1a1a] rounded-lg shadow-xl p-8 border border-[#333]">
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label
                htmlFor="dietary_restrictions"
                className="block text-sm font-medium mb-2 text-gray-300"
              >
                Dietary Restrictions
              </label>
              <input
                id="dietary_restrictions"
                type="text"
                value={dietaryRestrictions}
                onChange={(e) => setDietaryRestrictions(e.target.value)}
                className="w-full px-4 py-3 bg-[#2a2a2a] border border-[#444] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#646cff] focus:border-transparent transition-all"
                placeholder="e.g., vegetarian, vegan, gluten-free, dairy-free"
              />
              <p className="text-gray-500 text-xs mt-2">
                Separate multiple restrictions with commas
              </p>
            </div>

            <div>
              <label
                htmlFor="max_cooking_time"
                className="block text-sm font-medium mb-2 text-gray-300"
              >
                Maximum Cooking Time (minutes)
              </label>
              <input
                id="max_cooking_time"
                type="number"
                min="0"
                value={maxCookingTime || ""}
                onChange={(e) => setMaxCookingTime(parseInt(e.target.value) || 0)}
                className="w-full px-4 py-3 bg-[#2a2a2a] border border-[#444] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#646cff] focus:border-transparent transition-all"
                placeholder="e.g., 30"
              />
              <p className="text-gray-500 text-xs mt-2">
                Set to 0 for no time limit
              </p>
            </div>

            {error && (
              <div className="bg-red-500/10 border border-red-500/50 text-red-400 px-4 py-3 rounded-lg text-sm">
                {error}
              </div>
            )}

            {success && (
              <div className="bg-green-500/10 border border-green-500/50 text-green-400 px-4 py-3 rounded-lg text-sm">
                {success}
              </div>
            )}

            <button
              type="submit"
              disabled={saving}
              className="w-full bg-[#646cff] hover:bg-[#535bf2] disabled:bg-[#464a87] disabled:cursor-not-allowed text-white font-medium py-3 px-4 rounded-lg transition-colors duration-200"
            >
              {saving ? "Saving..." : "Save Preferences"}
            </button>
          </form>
        </div>

        {/* Info */}
        <div className="mt-6 bg-[#1a1a1a] rounded-lg p-6 border border-[#333]">
          <h2 className="text-lg font-semibold mb-3">How preferences work</h2>
          <ul className="space-y-2 text-gray-400 text-sm">
            <li className="flex items-start gap-2">
              <span className="text-[#646cff] mt-1">•</span>
              <span>
                Your dietary restrictions will be automatically included in all meal suggestions
              </span>
            </li>
            <li className="flex items-start gap-2">
              <span className="text-[#646cff] mt-1">•</span>
              <span>
                Recipes will be tailored to fit within your maximum cooking time
              </span>
            </li>
            <li className="flex items-start gap-2">
              <span className="text-[#646cff] mt-1">•</span>
              <span>
                You can update these preferences anytime
              </span>
            </li>
          </ul>
        </div>
      </div>
    </div>
  );
}

export default Preferences;
