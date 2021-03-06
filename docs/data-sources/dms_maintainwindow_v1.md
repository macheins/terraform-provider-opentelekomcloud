---
subcategory: "Distributed Message Service (DMS)"
---

# opentelekomcloud_dms_maintainwindow_v1

Use this data source to get the ID of an available OpenTelekomCloud DMS maintainwindow.

## Example Usage

```hcl
data "opentelekomcloud_dms_maintainwindow_v1" "maintainwindow1" {
  seq = 1
}
```

## Argument Reference

* `seq` - (Required) Indicates the sequential number of a maintenance time window.

* `begin` - (Optional) Indicates the time at which a maintenance time window starts.

* `end` - (Optional) Indicates the time at which a maintenance time window ends.

* `default` - (Optional) Indicates whether a maintenance time window is set to the default time segment.

## Attributes Reference

`id` is set to the ID of the found maintainwindow. In addition, the following attributes are exported:

* `seq` - See Argument Reference above.

* `begin` - See Argument Reference above.

* `end` - See Argument Reference above.

* `default` - See Argument Reference above.
