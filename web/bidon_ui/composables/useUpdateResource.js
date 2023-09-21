export default function ({
  path,
  message,
  onError = async () => {
    /* no operation function */
  },
}) {
  return useFormSubmit({
    path,
    message,
    method: "patch",
    onError: onError,
  });
}
