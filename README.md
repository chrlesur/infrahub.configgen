# ConfigGen - AI-Powered Infrastructure Configuration Generator

ConfigGen is a command-line tool that leverages AI.YOU to automatically generate YAML configurations for various types of IT infrastructure components, ensuring compliance with your JSON schema specifications.

## Features

- 🤖 AI-powered configuration generation
- 📋 JSON schema validation support
- 🔄 Streaming generation with real-time output
- 🎯 Multiple infrastructure types support (router, switch, firewall, server)
- 🔒 Secure authentication handling
- 📝 Clean YAML output
- 🎛️ Debug mode for troubleshooting

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/configgen.git

# Build the application
cd configgen
go build -o configgen
```

## Quick Start

The tool uses the official InfraHub schema, available at:
```
https://schema.infrahub.app/infrahub/schema/latest.json
```

Example usage with Claude assistant (ID: asst_dRvlfW2L6NkVyjeyvbU9tKBA):

```bash
configgen generate \
  --assistant "asst_dRvlfW2L6NkVyjeyvbU9tKBA" \
  --email "YOUREMAIL" \
  --password "YOURPASSWORD" \
  --type "router" \
  --context "Configuration d'un routeur Cisco 2600 pour un réseau d'entreprise avec 3 VLANs." \
  --schema docs/schema.json \
  --debug
```

This will generate a YAML configuration like:

```yaml
version: "1.0.0"
nodes:
- name: CiscoRouter2600
  namespace: Network
  description: "Configuration d'un routeur Cisco 2600 avec 3 VLANs"
  label: "Cisco 2600 Router"
  attributes:
  - name: hostname
    kind: string
    label: "Hostname"
    description: "Nom du routeur"
    min_length: 3
    max_length: 64
  - name: model
    kind: string
    label: "Model"
    default_value: "2600"
  - name: ios_version
    kind: string
    label: "IOS Version"
    description: "Version du système IOS"
  - name: mgmt_ip
    kind: string
    label: "Management IP"
    description: "Adresse IP de gestion"
  - name: vlan_config
    kind: list
    label: "VLAN Configuration"
    description: "Configuration des VLANs"
  - name: interfaces
    kind: list
    label: "Network Interfaces"
  - name: routing_protocol
    kind: dropdown
    label: "Routing Protocol"
    choices:
    - name: ospf
      label: "OSPF"
    - name: eigrp
      label: "EIGRP"
    - name: rip
      label: "RIP"
  - name: access_lists
    kind: list
    label: "Access Lists"
  - name: snmp_config
    kind: boolean
    label: "SNMP Enabled"
    default_value: true
  relationships:
  - name: connected_switches
    peer: NetworkSwitch
    label: "Connected Switches"
    cardinality: many
    description: "Switches connected to this router"
  - name: upstream_router
    peer: CiscoRouter2600
    label: "Upstream Router"
    cardinality: one
    optional: true
```

## Usage

### Basic Usage

```bash
./configgen generate \
  --email "your.email@example.com" \
  --password "your_password" \
  --assistant "asst_dRvlfW2L6NkVyjeyvbU9tKBA" \
  --schema "/path/to/schema.json" \
  --type "router" \
  --context "Configure a Cisco router with 3 VLANs for enterprise network"
```

### Quiet Mode (Output Only YAML)

```bash
./configgen generate \
  --quiet \
  --email "your.email@example.com" \
  --password "your_password" \
  --assistant "asst_dRvlfW2L6NkVyjeyvbU9tKBA" \
  --schema "/path/to/schema.json" \
  --type "router" \
  --context "Configure a Cisco router with 3 VLANs"
```

### Debug Mode

```bash
./configgen generate \
  --debug \
  --email "your.email@example.com" \
  --password "your_password" \
  --assistant "asst_dRvlfW2L6NkVyjeyvbU9tKBA" \
  --schema "/path/to/schema.json" \
  --type "router" \
  --context "Configure a Cisco router with 3 VLANs"
```

## Command Line Options

### Global Flags
- `--email`: Your AI.YOU email address
- `--password`: Your AI.YOU password
- `--assistant`: The ID of the AI.YOU assistant to use (default: asst_dRvlfW2L6NkVyjeyvbU9tKBA)
- `--config`: Path to configuration file (optional)
- `--debug`: Enable debug mode
- `--quiet`: Disable status messages

### Generate Command Flags
- `--type`: Type of machine (router, switch, firewall, server)
- `--context`: Natural language description of the desired configuration
- `--schema`: Path to JSON schema file

## Output

The tool generates a YAML file in the current directory with the naming format:
```
{type}_config_{timestamp}.yaml
```

Example:
```
router_config_20240122_153045.yaml
```

## Schema Validation

ConfigGen uses the official InfraHub JSON schema (https://schema.infrahub.app/infrahub/schema/latest.json) to validate and structure the generated configurations. The schema defines:

- Required fields and their types
- Allowed values and constraints
- Nested structures and relationships
- Validation rules

## Dependencies

- github.com/chrlesur/aiyou.golib
- github.com/spf13/cobra

## Error Handling

The tool includes comprehensive error handling for:
- Authentication failures
- Schema validation errors
- Network connectivity issues
- File system operations
- API response validation

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## Best Practices

1. Always validate the generated configuration before applying it
2. Use version control for your schema files
3. Start with the `--debug` flag when troubleshooting
4. Keep your context descriptions clear and specific
5. Use the `--quiet` flag for script integration

## License

This project is licensed under the GPL-3.0 License - see the [LICENSE](LICENSE) file for details.

## Support

For support, please:
1. Check the debug output (`--debug` flag)
2. Create an issue in the repository
3. Provide schema and context examples when reporting issues

## Roadmap

- [ ] Add configuration validation
- [ ] Support for multiple output formats
- [ ] Template system for common configurations
- [ ] Configuration diff tool
- [ ] Batch processing mode