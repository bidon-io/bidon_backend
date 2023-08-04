import { FilterMatchMode } from "primevue/api";

export const ResourceTableFields = {
  Id: { field: "id", header: "Id", sortable: true },
  App: {
    field: "appId",
    header: "App",
    link: {
      basePath: "/apps",
      dataField: "app",
      extractLinkData: ({ app }) => ({
        isValid: !!app,
        id: app?.id,
        linkText: `${app?.packageName} (${app?.platformId})`,
      }),
    },
    filter: {
      field: "appId",
      type: "select",
      matchMode: FilterMatchMode.EQUALS,
      placeholder: "App",
      extractOptions: (records) => [
        ...new Map(
          records.map(({ app }) => [
            app?.id,
            {
              label: `${app?.packageName} (${app?.platformId})`,
              value: app?.id,
            },
          ])
        ).values(),
      ],
    },
  },
  AccountType: { field: "accountType", header: "Account Type" },
  AdType: {
    field: "adType",
    header: "Ad Type",
    filter: {
      field: "adType",
      type: "select",
      matchMode: FilterMatchMode.EQUALS,
      placeholder: "AdType",
      extractOptions: (records) => [
        ...new Map(
          records.map(({ adType }) => [
            adType,
            { label: adType, value: adType },
          ])
        ).values(),
      ],
    },
  },
  BidFloor: { field: "bidFloor", header: "Bid Floor" },
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
    filter: {
      field: "accountId",
      type: "select",
      matchMode: FilterMatchMode.EQUALS,
      placeholder: "Account",
      extractOptions: (records) => [
        ...new Map(
          records.map(({ account }) => [
            account?.id,
            {
              label: `${account?.type?.split("::")[1]} (${account?.id})`,
              value: account?.id,
            },
          ])
        ).values(),
      ],
    },
  },
  HumanName: {
    field: "humanName",
    header: "Human Name",
    filter: {
      field: "humanName",
      type: "input",
      matchMode: FilterMatchMode.CONTAINS,
      placeholder: "Human Name",
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
