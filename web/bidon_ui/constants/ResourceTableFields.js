import { FilterMatchMode } from "primevue/api";
import { AdTypeEnum } from "~/types";

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
      type: "select-filter",
      matchMode: FilterMatchMode.EQUALS,
      placeholder: "App",
      loadOptions: async () => {
        const apps = await $apiFetch("/apps");
        return apps.map(({ id, packageName, platformId }) => ({
          label: `${packageName} (${platformId})`,
          value: id,
        }));
      },
      extractOptions: (records) => [
        ...new Map(
          records.map(({ app }) => [
            app?.id,
            {
              label: `${app?.packageName} (${app?.platformId})`,
              value: String(app?.id),
            },
          ]),
        ).values(),
      ],
    },
  },
  AccountType: { field: "accountType", header: "Account Type" },
  IsDefault: {
    field: "isDefault",
    header: "Default",
    filter: {
      field: "isDefault",
      type: "select",
      matchMode: FilterMatchMode.EQUALS,
      placeholder: "Default",
      loadOptions: async () => [
        {
          label: "True",
          value: "true",
        },
        {
          label: "False",
          value: "false",
        },
      ],
      extractOptions: () => [
        {
          label: "True",
          value: "true",
        },
        {
          label: "False",
          value: "false",
        },
      ],
    },
  },
  AdType: {
    field: "adType",
    header: "Ad Type",
    filter: {
      field: "adType",
      type: "select",
      matchMode: FilterMatchMode.EQUALS,
      placeholder: "AdType",
      loadOptions: async () =>
        Object.values(AdTypeEnum).map((adType) => ({
          label: adType,
          value: adType,
        })),
      extractOptions: (records) => [
        ...new Map(
          records.map(({ adType }) => [
            adType,
            { label: adType, value: adType },
          ]),
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
      loadOptions: async () => {
        const demandSources = await $apiFetch("/demand_sources");
        return demandSources.map(({ id, humanName }) => ({
          label: humanName,
          value: id,
        }));
      },
      extractOptions: (records) => [
        ...new Map(
          records.map(({ demandSource }) => [
            demandSource?.id,
            {
              label: demandSource?.humanName,
              value: String(demandSource?.id),
            },
          ]),
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
      type: "select-filter",
      matchMode: FilterMatchMode.EQUALS,
      placeholder: "Account",
      loadOptions: async () => {
        const accounts = await $apiFetch("/demand_source_accounts");
        return accounts.map(({ id, type, label }) => ({
          label: `(${type.split("::")[1]}) ${label ? label : `#${id}`}`,
          value: id,
        }));
      },
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
          ]),
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
  Name: {
    field: "name",
    header: "Name",
    filter: {
      field: "name",
      type: "input",
      matchMode: FilterMatchMode.CONTAINS,
      placeholder: "Name",
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
  AuctionKey: {
    field: "auctionKey",
    header: "Auction Key",
    copyable: true,
    filter: {
      field: "auctionKey",
      type: "input",
      matchMode: FilterMatchMode.CONTAINS,
      placeholder: "Auction Key",
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
      type: "select-filter",
      matchMode: FilterMatchMode.EQUALS,
      placeholder: "Segment",
      loadOptions: async () => {
        const segments = await $apiFetch("/segments");
        return segments.map(({ id, name }) => ({
          label: name,
          value: id,
        }));
      },
      extractOptions: (records) => [
        ...new Map(
          records.map(({ segment }) => [
            segment?.id,
            {
              label: segment?.name,
              value: String(segment?.id),
            },
          ]),
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
      type: "select-filter",
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
          ]),
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
