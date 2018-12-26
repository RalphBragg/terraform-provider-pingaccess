package pingaccess

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/iwarapter/pingaccess-sdk-go/pingaccess"
)

func resourcePingAccessApplicationResource() *schema.Resource {
	return &schema.Resource{
		Create: resourcePingAccessApplicationResourceCreate,
		Read:   resourcePingAccessApplicationResourceRead,
		Update: resourcePingAccessApplicationResourceUpdate,
		Delete: resourcePingAccessApplicationResourceDelete,

		Schema: map[string]*schema.Schema{
			anonymous: &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			applicationID: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			auditLevel: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			defaultAuthTypeOverride: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			enabled: &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			methods: &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			name: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			pathPrefixes: &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			policy: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			rootResource: &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourcePingAccessApplicationResourceCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] resourcePingAccessApplicationResourceCreate")
	anonymous := d.Get(anonymous).(bool)
	application_id := d.Get(applicationID).(string)
	audit_level := d.Get(auditLevel).(string)
	default_auth_type_override := d.Get(defaultAuthTypeOverride).(string)
	enabled := d.Get(enabled).(bool)
	// methods := d.Get("methods").([]*string)
	methods := expandStringList(d.Get(methods).(*schema.Set).List())
	name := d.Get(name).(string)
	// path_prefixes := d.Get("path_prefixes").([]*string)
	path_prefixes := expandStringList(d.Get(pathPrefixes).(*schema.Set).List())
	// policy := d.Get("policy").(map[string]interface{})
	root_resource := d.Get(rootResource).(bool)

	appId, _ := strconv.Atoi(application_id)

	input := pingaccess.AddApplicationResourceCommandInput{
		Id: application_id,
		Body: pingaccess.ResourceView{
			Anonymous:               anonymous,
			ApplicationId:           appId,
			AuditLevel:              audit_level,
			DefaultAuthTypeOverride: default_auth_type_override,
			Enabled:                 enabled,
			Methods:                 methods,
			Name:                    name,
			PathPrefixes:            path_prefixes,
			// Policy:                  policy,
			RootResource: root_resource,
		},
	}
	svc := m.(*pingaccess.Client).Applications

	result, _, err := svc.AddApplicationResourceCommand(&input)
	if err != nil {
		return fmt.Errorf("Error creating application: %s", err)
	}

	d.SetId(strconv.Itoa(result.Id))
	return resourcePingAccessApplicationResourceReadResult(d, &input.Body)
}

func resourcePingAccessApplicationResourceRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] resourcePingAccessApplicationResourceRead")
	svc := m.(*pingaccess.Client).Applications

	input := &pingaccess.GetApplicationResourceCommandInput{
		ApplicationId: d.Get("application_id").(string),
		ResourceId:    d.Id(),
	}

	log.Printf("[INFO] ResourceID: %s", d.Id())
	log.Printf("[INFO] GetApplicationResourceCommandInput: %s", input.ApplicationId)
	result, _, _ := svc.GetApplicationResourceCommand(input)
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(result)
	rs := pingaccess.ResourceView{}
	json.NewDecoder(b).Decode(&rs)

	return resourcePingAccessApplicationResourceReadResult(d, &rs)
}

func resourcePingAccessApplicationResourceUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] resourcePingAccessApplicationResourceUpdate")
	anonymous := d.Get(anonymous).(bool)
	application_id := d.Get(applicationID).(string)
	audit_level := d.Get(auditLevel).(string)
	default_auth_type_override := d.Get(defaultAuthTypeOverride).(string)
	enabled := d.Get(enabled).(bool)
	// methods := d.Get("methods").([]*string)
	methods := expandStringList(d.Get(methods).(*schema.Set).List())
	name := d.Get(name).(string)
	// path_prefixes := d.Get("path_prefixes").([]*string)
	path_prefixes := expandStringList(d.Get(pathPrefixes).(*schema.Set).List())
	// policy := d.Get("policy").(map[string]interface{})
	root_resource := d.Get(rootResource).(bool)

	id, _ := strconv.Atoi(d.Id())
	appId, _ := strconv.Atoi(application_id)

	input := pingaccess.UpdateApplicationResourceCommandInput{
		ApplicationId: application_id,
		ResourceId:    d.Id(),
		Body: pingaccess.ResourceView{
			Anonymous:               anonymous,
			ApplicationId:           appId,
			AuditLevel:              audit_level,
			DefaultAuthTypeOverride: default_auth_type_override,
			Enabled:                 enabled,
			Id:                      id,
			Methods:                 methods,
			Name:                    name,
			PathPrefixes:            path_prefixes,
			// Policy:                  policy,
			RootResource: root_resource,
		},
	}
	input.ApplicationId = d.Get("application_id").(string)

	svc := m.(*pingaccess.Client).Applications

	_, _, err := svc.UpdateApplicationResourceCommand(&input)
	if err != nil {
		return fmt.Errorf("Error updating application: %s", err)
	}
	log.Println("[DEBUG] End - resourcePingAccessApplicationResourceUpdate")
	return nil
}

func resourcePingAccessApplicationResourceDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] resourcePingAccessApplicationResourceDelete")
	return nil
}

func resourcePingAccessApplicationResourceReadResult(d *schema.ResourceData, rv *pingaccess.ResourceView) error {
	log.Printf("[INFO] resourcePingAccessApplicationResourceReadResult")
	if err := d.Set(anonymous, rv.Anonymous); err != nil {
		return err
	}
	// anonymous := d.Get("anonymous").(bool)
	if err := d.Set(applicationID, strconv.Itoa(rv.ApplicationId)); err != nil {
		return err
	}
	// application_id := d.Get("application_id").(string)
	if err := d.Set(auditLevel, rv.AuditLevel); err != nil {
		return err
	}
	// audit_level := d.Get("audit_level").(string)
	if err := d.Set(defaultAuthTypeOverride, rv.DefaultAuthTypeOverride); err != nil {
		return err
	}
	// default_auth_type_override := d.Get("default_auth_type_override").(string)
	if err := d.Set(enabled, rv.Enabled); err != nil {
		return err
	}
	// enabled := d.Get("enabled").(bool)
	// methods := d.Get("methods").([]*string)
	if err := d.Set(methods, rv.Methods); err != nil {
		return err
	}
	// methods := expandStringList(d.Get("methods").(*schema.Set).List())
	if err := d.Set(name, rv.Name); err != nil {
		return err
	}
	// name := d.Get("name").(string)
	// path_prefixes := d.Get("path_prefixes").([]*string)
	if err := d.Set(pathPrefixes, rv.PathPrefixes); err != nil {
		return err
	}
	// path_prefixes := expandStringList(d.Get("path_prefixes").(*schema.Set).List())
	// if err := d.Set("policy", rv.Policy); err != nil {
	// 	return err
	// }
	// policy := d.Get("policy").(map[string]interface{})
	if err := d.Set(rootResource, rv.RootResource); err != nil {
		return err
	}
	return nil
}