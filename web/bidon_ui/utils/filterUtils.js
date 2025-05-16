import { decamelize } from "humps";

/**
 * Transforms filter values for special cases like AdTypeWithFormat
 * @param {string} field - The filter field name
 * @param {any} value - The filter value
 * @returns {Object} - Object with transformed field name and value
 */
export const transformFilterValue = (field, value) => {
  if (field === "adTypeWithFormat" && value) {
    if (typeof value === "object" && value.adType) {
      return {
        fields: {
          ad_type: value.adType,
          format: value.format || "",
        },
      };
    }
  }

  return {
    fields: {
      [decamelize(field)]: value,
    },
  };
};

/**
 * Builds query parameters from filters
 * @param {Object} filters - The filters object
 * @param {number} page - Current page
 * @param {number} limit - Items per page
 * @returns {Object} - Query parameters object
 */
export const buildQueryParams = (filters, page, limit) => {
  const query = {
    page,
    limit,
  };

  Object.entries(filters)
    .filter(([, filter]) => filter?.value)
    .forEach(([key, filter]) => {
      const { fields } = transformFilterValue(key, filter.value);

      // Add all transformed fields to the query
      Object.entries(fields).forEach(([fieldName, fieldValue]) => {
        query[fieldName] = fieldValue;
      });
    });

  return query;
};
