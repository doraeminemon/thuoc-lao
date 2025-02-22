// stackbit.config.ts
import { defineStackbitConfig } from "@stackbit/types";
import { GitContentSource } from "@stackbit/cms-git";

export default defineStackbitConfig({
  stackbitVersion: "~0.6.0",
  contentSources: [
    new GitContentSource({
      rootPath: __dirname,
      contentDirs: ["content"],
      models: [
        {
          name: "Page",
          type: "page",
          urlPath: "/{slug}",
          filePath: "content/pages/{slug}.json",
          fields: [{ name: "title", type: "string", required: true }],
        },
      ],
      assetsConfig: {
        referenceType: "static",
        staticDir: "public",
        uploadDir: "images",
        publicPath: "/",
      },
    }),
  ],
  siteMap: ({ documents, models }) => {
    // 1. Filter all page models
    const pageModels = models.filter((m) => m.type === "page");
    return (
      documents
        // 2. Filter all documents which are of a page model
        .filter((d) => pageModels.some((m) => m.name === d.modelName))
        // 3. Map each document to a SiteMapEntry
        .map((document) => {
          // Map the model name to its corresponding URL
          const urlModel = (() => {
            switch (document.modelName) {
              case "Page":
                return "otherPage";
              case "Blog":
                return "otherBlog";
              default:
                return null;
            }
          })();

          return {
            stableId: document.id,
            urlPath: `/${urlModel}/${document.id}`,
            document,
            isHomePage: false,
          };
        })
        .filter(Boolean)
    );
  },
});
