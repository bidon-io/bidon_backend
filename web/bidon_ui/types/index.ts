export { ResourceLink, AdTypeEnum, AdType };

declare global {
  type SomeType = [boolean, string, number];

  interface ResourceLink {
    basePath: string;
    extractLinkData: (data: any) => {
      linkText: string;
      id: number;
      isValid: boolean;
    };
  }

  interface AssociatedResourcesLink {
    extractLinkData: (data: any) => {
      label: string;
      path: string;
    };
  }

  interface ResourcePermissions {
    create?: boolean;
    read?: boolean;
    update?: boolean;
    delete?: boolean;
  }
}

enum AdTypeEnum {
  Banner = "banner",
  Interstitial = "interstitial",
  Rewarded = "rewarded",
}

type AdType = AdTypeEnum.Banner | AdTypeEnum.Interstitial | AdTypeEnum.Rewarded;
