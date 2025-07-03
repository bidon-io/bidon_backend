import { ref, computed } from "vue";
import { useAsyncData } from "#app";
import axios from "@/services/ApiService";

export function useAppDemandProfileValidation(appId) {
  const warnings = ref(new Map());

  // Fetch demand sources to map api_key to demand_source_id
  const { data: demandSources } = useAsyncData(
    "demand-sources-for-validation",
    async () => {
      try {
        const response = await axios.get("/demand_sources");
        return response.data;
      } catch (error) {
        console.error("Failed to fetch demand sources:", error);
        return [];
      }
    },
    {
      default: () => [],
      server: false,
    },
  );

  // Create a map from api_key to demand_source_id
  const demandSourceMap = computed(() => {
    if (!demandSources.value) return new Map();

    const map = new Map();
    demandSources.value.forEach((source) => {
      map.set(source.apiKey, source.id);
    });
    return map;
  });

  // Fetch app demand profiles for the current app
  const { data: appDemandProfiles, refresh: refreshProfiles } = useAsyncData(
    `app-demand-profiles-${appId.value}`,
    async () => {
      if (!appId.value) return [];

      try {
        const response = await axios.get("/app_demand_profiles_collection", {
          params: {
            app_id: appId.value,
            limit: 1000, // Get all profiles for this app
          },
        });
        return response.data?.items || [];
      } catch (error) {
        console.error("Failed to fetch app demand profiles:", error);
        return [];
      }
    },
    {
      default: () => [],
      server: false,
      watch: [appId],
    },
  );

  // Create a map from demand_source_id to app demand profile
  const profileMap = computed(() => {
    if (!appDemandProfiles.value) return new Map();

    const map = new Map();
    appDemandProfiles.value.forEach((profile) => {
      map.set(profile.demandSourceId, profile);
    });
    return map;
  });

  const validateNetworkEnabled = async (networkApiKey) => {
    // Get demand_source_id from api_key
    const demandSourceId = demandSourceMap.value.get(networkApiKey);

    if (!demandSourceId) {
      console.warn(`No demand source found for api_key: ${networkApiKey}`);
      warnings.value.set(networkApiKey, {
        message: `No demand source found for network "${networkApiKey}". Please check if the demand source exists.`,
        severity: "warn",
      });
      return;
    }

    // Get the corresponding app demand profile
    const profile = profileMap.value.get(demandSourceId);

    if (!profile) {
      // No profile exists - this is a warning condition
      warnings.value.set(networkApiKey, {
        message: `No App Demand Profile found for network "${networkApiKey}". Please create and enable an App Demand Profile for this network.`,
        severity: "warn",
      });
      return;
    }

    if (!profile.enabled) {
      // Profile exists but is disabled - this is a warning condition
      warnings.value.set(networkApiKey, {
        message: `App Demand Profile for network "${networkApiKey}" is disabled. Please enable the corresponding App Demand Profile.`,
        severity: "warn",
      });
      return;
    }

    // Profile exists and is enabled - clear any existing warning
    warnings.value.delete(networkApiKey);
  };

  const validateNetworkDisabled = (networkApiKey) => {
    // Clear warning when network is disabled
    warnings.value.delete(networkApiKey);
  };

  const clearAllWarnings = () => {
    warnings.value.clear();
  };

  const warningMessages = computed(() => {
    return Array.from(warnings.value.values());
  });

  const hasWarnings = computed(() => {
    return warnings.value.size > 0;
  });

  return {
    validateNetworkEnabled,
    validateNetworkDisabled,
    clearAllWarnings,
    warningMessages,
    hasWarnings,
    refreshProfiles,
  };
}
