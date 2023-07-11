export const ResourceTableFields = {
  Id: { field: "id", header: "Id", sortable: true },
  App: {
    field: "appId",
    header: "App",
    link: {
      basePath: "/apps",
      textField: "packageName",
      dataField: "app",
      extractLinkData: ({ app }) => ({
        isValid: !!app,
        id: app?.id,
        linkText: `${app?.packageName} (${app?.platformId})`,
      }),
    },
  },
  DemandSource: {
    field: "demandSourceId",
    header: "Demand Source",
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
    field: "accountId",
    header: "Account",
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
    field: "segmentId",
    header: "Segment",
    link: {
      basePath: "/segments",
      extractLinkData: ({ segment }) => ({
        isValid: !!segment,
        id: segment?.id,
        linkText: segment?.name,
      }),
    },
  },
  User: {
    field: "userId",
    header: "User",
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
