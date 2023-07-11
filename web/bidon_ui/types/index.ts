export { ResourceLink };

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
}
