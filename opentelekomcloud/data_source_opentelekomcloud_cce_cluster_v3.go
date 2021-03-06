package opentelekomcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/cce/v3/clusters"
)

func dataSourceCCEClusterV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCCEClusterV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"billing_mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"highway_subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"container_network_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"container_network_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"internal": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_otc": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_authority_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"certificate_users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_certificate_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_key_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCCEClusterV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.cceV3Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Unable to create opentelekomcloud CCE client : %s", err)
	}

	listOpts := clusters.ListOpts{
		ID:    d.Id(),
		Name:  d.Get("name").(string),
		Type:  d.Get("cluster_type").(string),
		Phase: d.Get("status").(string),
		VpcID: d.Get("vpc_id").(string),
	}

	refinedClusters, err := clusters.List(cceClient, listOpts)
	log.Printf("[DEBUG] Value of allClusters: %#v", refinedClusters)
	if err != nil {
		return fmt.Errorf("Unable to retrieve clusters: %s", err)
	}

	if len(refinedClusters) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedClusters) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	cluster := refinedClusters[0]

	log.Printf("[DEBUG] Retrieved Clusters using given filter %s: %+v", cluster.Metadata.Id, cluster)

	d.SetId(cluster.Metadata.Id)

	d.Set("name", cluster.Metadata.Name)
	d.Set("flavor_id", cluster.Spec.Flavor)
	d.Set("description", cluster.Spec.Description)
	d.Set("cluster_version", cluster.Spec.Version)
	d.Set("cluster_type", cluster.Spec.Type)
	d.Set("billing_mode", cluster.Spec.BillingMode)
	d.Set("vpc_id", cluster.Spec.HostNetwork.VpcId)
	d.Set("subnet_id", cluster.Spec.HostNetwork.SubnetId)
	d.Set("highway_subnet_id", cluster.Spec.HostNetwork.HighwaySubnet)
	d.Set("container_network_cidr", cluster.Spec.ContainerNetwork.Cidr)
	d.Set("container_network_type", cluster.Spec.ContainerNetwork.Mode)
	d.Set("status", cluster.Status.Phase)
	d.Set("internal", cluster.Status.Endpoints[0].Internal)
	d.Set("external", cluster.Status.Endpoints[0].External)
	d.Set("external_otc", cluster.Status.Endpoints[0].ExternalOTC)
	d.Set("region", GetRegion(d, config))

	cert, err := clusters.GetCert(cceClient, d.Id()).Extract()
	if err != nil {
		log.Printf("Error retrieving opentelekomcloud CCE cluster cert: %s", err)
	}

	// Set Certificate Clusters
	var clusterList []map[string]interface{}
	for _, clusterObj := range cert.Clusters {
		clusterCert := make(map[string]interface{})
		clusterCert["name"] = clusterObj.Name
		clusterCert["server"] = clusterObj.Cluster.Server
		clusterCert["certificate_authority_data"] = clusterObj.Cluster.CertAuthorityData
		clusterList = append(clusterList, clusterCert)
	}
	d.Set("certificate_clusters", clusterList)

	// Set Certificate Users
	var userList []map[string]interface{}
	for _, userObj := range cert.Users {
		userCert := make(map[string]interface{})
		userCert["name"] = userObj.Name
		userCert["client_certificate_data"] = userObj.User.ClientCertData
		userCert["client_key_data"] = userObj.User.ClientKeyData
		userList = append(userList, userCert)
	}
	d.Set("certificate_users", userList)

	return nil
}
