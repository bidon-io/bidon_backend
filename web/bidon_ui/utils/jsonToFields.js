/**
 * Formats a camelCase key into a readable label with spaces and capitalization
 * @param {string} key - The camelCase key to format
 * @returns {string} - Formatted label with spaces and capitalization
 */
export const formatLabel = (key) => {
  return key
    .replace(/([A-Z])/g, " $1") // Add space before capital letters
    .replace(/^./, (str) => str.toUpperCase()); // Capitalize first letter
};

/**
 * Converts a JSON object to an array of field objects for ResourceCard
 * @param {Object} jsonData - The JSON data to convert
 * @param {string} prefix - Optional prefix for the key (e.g., 'data')
 * @param {string} type - The type of field to create (default: 'static')
 * @param {boolean} copyable - Whether the field is copyable (default: false)
 * @returns {Array} - Array of field objects
 */
export const jsonToFields = (
  jsonData,
  prefix = "",
  type = "",
  copyable = false,
) => {
  if (!jsonData || typeof jsonData !== "object") {
    return [];
  }

  return Object.keys(jsonData).map((key) => {
    const fieldKey = prefix ? `${prefix}.${key}` : key;

    return {
      label: formatLabel(key),
      key: fieldKey,
      value: jsonData[key],
      type: type || undefined,
      copyable,
    };
  });
};
