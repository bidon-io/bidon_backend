import { FilterMatchMode } from "primevue/api";

export const ResourceTableFields = {
  Id: { field: "id", header: "Id", sortable: true },
  PublicUid: { field: "publicUid", header: "Public UID" },
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
              value: String(app?.id),
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
    filter: {
      field: "demandSourceId",
      type: "select",
      matchMode: FilterMatchMode.EQUALS,
      placeholder: "Demand Source",
      extractOptions: (records) => [
        ...new Map(
          records.map(({ demandSource }) => [
            demandSource?.id,
            {
              label: demandSource?.humanName,
              value: String(demandSource?.id),
            },
          ])
        ).values(),
      ],
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
              label: `(${account?.type?.split("::")[1]}) ${
                account?.label ? account?.label : `#${account?.id}`
              }`,
              value: String(account?.id),
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
  AppName: {
    field: "humanName",
    header: "App Name",
    filter: {
      field: "humanName",
      type: "input",
      matchMode: FilterMatchMode.CONTAINS,
      placeholder: "App Name",
    },
  },
  Label: {
    field: "label",
    header: "Label",
    filter: {
      field: "label",
      type: "input",
      matchMode: FilterMatchMode.CONTAINS,
      placeholder: "Label",
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
    filter: {
      field: "segmentId",
      type: "select",
      matchMode: FilterMatchMode.EQUALS,
      placeholder: "Segment",
      extractOptions: (records) => [
        ...new Map(
          records.map(({ segment }) => [
            segment?.id,
            {
              label: segment?.name,
              value: String(segment?.id),
            },
          ])
        ).values(),
      ],
    },
  },
  Owner: {
    field: "userId",
    header: "Owner",
    link: {
      basePath: "/users",
      extractLinkData: ({ user }) => ({
        isValid: !!user,
        id: user?.id,
        linkText: user?.email,
      }),
    },
    filter: {
      field: "userId",
      type: "select",
      matchMode: FilterMatchMode.EQUALS,
      placeholder: "Owner",
      extractOptions: (records) => [
        ...new Map(
          records.map(({ user }) => [
            user?.id,
            {
              label: user?.email,
              value: String(user?.id),
            },
          ])
        ).values(),
      ],
    },
  },
  Email: {
    field: "email",
    header: "Email",
    filter: {
      field: "email",
      type: "input",
      matchMode: FilterMatchMode.CONTAINS,
      placeholder: "Email",
    },
  },
};
