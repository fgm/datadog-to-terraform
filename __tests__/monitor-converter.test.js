import { generateTerraformCode } from "../src/lib/monitor-converter";
import monitorJSON from "../examples/monitor.json";

it("converts correctly", () => {
  expect(generateTerraformCode("monitor_1", monitorJSON)).toMatchSnapshot();
});
