export const ResourceCardFields = {
  Id: { label: "ID", key: "id" },
  App: {
    label: "App",
    key: "appId",
    type: "link",
    link: {
      basePath: "/apps",
      extractLinkData: ({ app }) => ({
        isValid: !!app,
        id: app?.id,
        linkText: `${app?.packageName} (${app?.platformId})`,
      }),
    },
  },
  DemandSource: {
    label: "Demand Source",
    key: "demandSourceId",
    type: "link",
    link: {
      basePath: "/demand_sources",
      extractLinkData: ({ demandSource }) => ({
        isValid: !!demandSource,
        id: demandSource?.id,
        linkText: demandSource?.humanName,
      }),
    },
  },
  DemandSourceAccount: {
    label: "Account",
    key: "accountId",
    type: "link",
    link: {
      basePath: "/demand_source_accounts",
      extractLinkData: ({ account }) => ({
        isValid: !!account,
        id: account?.id,
        linkText: `${account?.type?.split("::")[1]} (${account?.id})`,
      }),
    },
  },
  Segment: {
    label: "Segment",
    key: "segmentId",
    type: "link",
    link: {
      basePath: "segments",
      extractLinkData: ({ segment }) => ({
        isValid: !!segment,
        id: segment?.id,
        linkText: segment?.name,
      }),
    },
  },
  User: {
    label: "User",
    key: "userId",
    type: "link",
    link: {
      basePath: "/users",
      extractLinkData: ({ user }) => ({
        isValid: !!user,
        id: user?.id,
        linkText: user?.email,
      }),
    },
  },
};
