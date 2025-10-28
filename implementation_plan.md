# Wild Workouts Template-Based DDD Implementation Plan

## Overview

This document outlines our plan to refactor the Wild Workouts codebase to extract the essential complexity from the accidental complexity, using a template-based approach with Protocol Buffers as the schema definition language.

## Important: Exploratory Nature of This Work

**This is a design exploration and feasibility study, not a production implementation.**

The purpose of this exercise is to:
- Explore whether template-based DDD code generation is viable
- Test the ergonomics of a contract-driven development approach
- Validate whether we can cleanly separate essential from accidental complexity
- Identify patterns that emerge when applying templates to real-world DDD code
- Create a proof of concept that demonstrates the idea

We're building a tangible demonstration to understand what this approach would look like in practice. The generated code may not be perfect, and we may discover limitations or issues along the way. That's expected and valuable - the goal is to learn whether this architectural approach has merit.

## Goals

1. Separate essential business logic from accidental infrastructure complexity
2. Enable code generation for repetitive patterns
3. Maintain the DDD principles while reducing boilerplate
4. Create a more maintainable and evolvable codebase

## Research on Protocol Buffer Editions

Based on our research, Protocol Buffer editions provide a way to evolve the protobuf language incrementally:

- Editions replace the proto2/proto3 syntax with a versioned approach (e.g., `edition = "2024"`)
- Features can be set at different scopes (file, message, field) with specific default values for each edition
- Custom options are implemented through extensions to existing option messages
- Custom options require importing specific option definitions
- Protocol buffer definitions can serve as contracts for required business logic implementations

For our DDD implementation, we'll use Protocol Buffer edition 2024 (the latest stable version) and extend it with custom options for our domain concepts. Our approach will follow a contract-driven development pattern, where proto definitions specify what business logic must be implemented.

## File Structure

```
/wild-workouts-go-ddd-example/
  /proto/                           # Domain definitions in proto format
    /options/                       # Custom option definitions
      ddd_options.proto             # DDD-specific options (entity, value object, etc.)
    /domain/                        # Core domain entities
      training.proto                # Training entity with behavior specifications
      user.proto                    # User definitions
    /application/                   # Application operations
      commands.proto                # All commands with their inputs/behaviors
      queries.proto                 # All queries with their inputs/outputs
    
  /templates/                       # Code generation templates
    /domain/                        # Domain layer templates
      entity.tpl                    # Entity struct and methods
      repository.tpl                # Repository interface
    /application/                   # Application layer templates
      command_handler.tpl           # Command handler template
      query_handler.tpl             # Query handler template
    /adapters/                      # Adapter layer templates
      firestore_repo.tpl            # Firestore implementation
      http_handler.tpl              # HTTP handler
    /contracts/                     # Business logic contract templates
      domain_method.tpl             # Domain method implementation contract
      command_logic.tpl             # Command business logic contract

  /business_logic/                  # Non-generated unique business logic
    /training/                      # Core business rules that can't be templated
      cancel.go                     # Cancel training business logic
      reschedule.go                 # Reschedule training business logic
    /stubs/                         # Generated contract stubs (for reference)
      cancel_stub.go                # Generated stub for cancel logic
      reschedule_stub.go            # Generated stub for reschedule logic
```

## Implementation Steps

### 1. Define Custom DDD Options

Create `ddd_options.proto` with custom options for DDD concepts by extending the standard Protocol Buffer options:

```proto
edition = "2024";
package wildworkouts.options;

import "google/protobuf/descriptor.proto";

extend google.protobuf.MessageOptions {
  bool entity = 50001;
  bool value_object = 50002;
  bool aggregate_root = 50003;
}

extend google.protobuf.FieldOptions {
  bool identifier = 60001;
  bool required = 60002;
  bool invariant = 60003; // Field with validation rules
}

extend google.protobuf.MethodOptions {
  bool domain_service = 70001;
  bool command = 70002;
  bool query = 70003;
  string validation = 70004; // Validation rules
  string authorization = 70005; // Auth rules
}

// Step type definitions
enum StepType {
  STEP_TYPE_UNSPECIFIED = 0;
  STEP_TYPE_LOAD = 1;         // Load entity from repository
  STEP_TYPE_SAVE = 2;         // Save entity to repository
  STEP_TYPE_CALL_METHOD = 3;  // Call domain method
  STEP_TYPE_CALL_SERVICE = 4; // Call external service
  STEP_TYPE_UPDATE = 5;       // Update an entity
  STEP_TYPE_CONDITION = 6;    // Conditional logic
}

// Defines a single execution step
message ExecutionStep {
  StepType type = 1;
  string target = 2;          // Entity or service being operated on
  string operation = 3;       // Method or operation to call
  repeated string parameters = 4;  // Parameters or conditions
  string condition = 5;       // Optional condition for conditional steps
  string comment = 6;         // Optional human-readable description
}

// Defines the entire execution flow
message ExecutionFlow {
  repeated ExecutionStep steps = 1;
}

// Additional options for application-specific concepts
extend google.protobuf.MessageOptions {
  repeated string dependencies = 80001; // Service dependencies
  ExecutionFlow execution_flow = 80002;  // Type-safe execution flow
  bool requires_implementation = 80003; // Whether this message requires business logic implementation
}

// Domain method definition (since RPCs can only be in services, we use a message-based approach)
message DomainMethod {
  string name = 1;              // Method name
  string request_type = 2;      // Request message type (optional)
  string response_type = 3;     // Response message type (optional)
  string validation = 4;        // Validation rule name
  string comment = 5;           // Human-readable description
}

extend google.protobuf.MessageOptions {
  repeated DomainMethod domain_methods = 90001;  // Domain methods available on this entity
}
```

### 2. Define Domain Entities in Proto

For each domain entity:
1. Define the entity structure using Protocol Buffers
2. Apply the custom DDD options to mark domain concepts
3. Define domain behaviors (methods)
4. Add validation rules

Example: `training.proto`

```proto
edition = "2024";
package wildworkouts.domain;

import "google/protobuf/timestamp.proto";
import "options/ddd_options.proto";

message Training {
  option (wildworkouts.options.entity) = true;
  option (wildworkouts.options.aggregate_root) = true;
  
  // Define domain methods available on this entity
  option (wildworkouts.options.domain_methods) = {
    name: "Cancel",
    validation: "cannot_cancel_if_too_late",
    comment: "Cancels the training session"
  };
  option (wildworkouts.options.domain_methods) = {
    name: "Reschedule",
    request_type: "RescheduleRequest",
    validation: "cannot_reschedule_if_too_late",
    comment: "Reschedules the training to a new time"
  };
  
  string uuid = 1 [(wildworkouts.options.identifier) = true];
  string user_uuid = 2 [(wildworkouts.options.required) = true];
  string user_name = 3 [(wildworkouts.options.required) = true];
  
  google.protobuf.Timestamp time = 4 [(wildworkouts.options.required) = true];
  string notes = 5;
  
  google.protobuf.Timestamp proposed_time = 6;
  UserType move_proposed_by = 7;
  
  bool canceled = 8;
  
  // Supporting message types for domain methods
  message RescheduleRequest {
    google.protobuf.Timestamp new_time = 1;
  }
}

enum UserType {
  USER_TYPE_UNSPECIFIED = 0;
  USER_TYPE_ATTENDEE = 1;
  USER_TYPE_TRAINER = 2;
}

message User {
  option (wildworkouts.options.value_object) = true;
  
  string uuid = 1 [(wildworkouts.options.required) = true];
  UserType type = 2 [(wildworkouts.options.required) = true];
}
```

### 3. Define Application Layer in Proto

1. Define commands with:
   - Input parameters
   - Dependencies (repositories, services)
   - Execution flow steps

2. Define queries with:
   - Input parameters
   - Return structures
   - Dependencies

Example: `commands.proto`

```proto
edition = "2024";
package wildworkouts.application;

import "google/protobuf/timestamp.proto";
import "domain/training.proto";
import "options/ddd_options.proto";

// Command definitions
message CancelTrainingCommand {
  option (wildworkouts.options.command) = true;
  option (wildworkouts.options.requires_implementation) = true;
  
  // Dependencies can be set as a repeated field or multiple options
  // Using the compact form here - actual syntax depends on protoc version
  option (wildworkouts.options.dependencies) = "TrainingRepository";
  option (wildworkouts.options.dependencies) = "UserService";
  option (wildworkouts.options.dependencies) = "TrainerService";
  
  option (wildworkouts.options.execution_flow) = {
    steps: [
      {
        type: STEP_TYPE_LOAD,
        target: "Training",
        operation: "GetTraining",
        parameters: ["cmd.TrainingUUID", "cmd.User"],
        comment: "Load training from repository"
      },
      {
        type: STEP_TYPE_CALL_METHOD,
        target: "training",
        operation: "Cancel",
        comment: "Call training.Cancel()"
      },
      {
        type: STEP_TYPE_CONDITION,
        condition: "balanceDelta := training.CancelBalanceDelta(cmd.User.Type()); balanceDelta != 0",
        comment: "Check if balance needs updating"
      },
      {
        type: STEP_TYPE_CALL_SERVICE,
        target: "userService",
        operation: "UpdateTrainingBalance",
        parameters: ["training.UserUUID()", "balanceDelta"],
        condition: "balanceDelta != 0",
        comment: "Update user balance if needed"
      },
      {
        type: STEP_TYPE_CALL_SERVICE,
        target: "trainerService",
        operation: "CancelTraining",
        parameters: ["training.Time()"],
        comment: "Update trainer availability"
      },
      {
        type: STEP_TYPE_SAVE,
        target: "training",
        operation: "UpdateTraining",
        comment: "Save training to repository"
      }
    ]
  };
  
  string training_uuid = 1 [(wildworkouts.options.required) = true];
  wildworkouts.domain.User user = 2 [(wildworkouts.options.required) = true];
}
```

### 4. Create Templates for Code Generation

Develop templates that process the proto definitions and generate code following the patterns in Wild Workouts.

**Template Engine**: We'll use the template system examined at the beginning of this conversation (based on `buf.build/go/template`), which provides:
- Proto parsing with custom option support
- Access to message metadata and options
- Template functions for code generation
- Integration with the protobuf compilation pipeline

Example: Entity template (`entity.tpl`)

```
package {{ .Package }}

import (
	"time"
	"errors"
)

// {{ .Message.Name }} is a domain entity
{{- if .Message.IsAggregateRoot }}
// This is an Aggregate Root
{{- end }}
type {{ .Message.Name }} struct {
	{{- range .Message.Fields }}
	{{ .Name | camelCase }} {{ .Type | goType }}
	{{- end }}
}

// New{{ .Message.Name }} creates a new instance of {{ .Message.Name }}
func New{{ .Message.Name }}(
	{{- range .Message.RequiredFields }}
	{{ .Name | camelCase }} {{ .Type | goType }},
	{{- end }}
) (*{{ .Message.Name }}, error) {
	// Validate required fields
	{{- range .Message.RequiredFields }}
	{{- if eq .Type "string" }}
	if {{ .Name | camelCase }} == "" {
		return nil, errors.New("empty {{ .Name }}")
	}
	{{- end }}
	{{- if eq .Type "google.protobuf.Timestamp" }}
	if {{ .Name | camelCase }}.IsZero() {
		return nil, errors.New("zero {{ .Name }}")
	}
	{{- end }}
	{{- end }}
	
	return &{{ .Message.Name }}{
		{{- range .Message.RequiredFields }}
		{{ .Name | camelCase }}: {{ .Name | camelCase }},
		{{- end }}
	}, nil
}

// Import the business logic package
import (
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/business_logic/{{ .Package | lowercase }}"
)

{{ range .Message.DomainMethods }}
// {{ .Name }} implements domain behavior
// Delegates to business logic implementation in business_logic/{{ $.Package | lowercase }}/{{ .Name | lowercase }}.go
func (e *{{ $.Message.Name }}) {{ .Name }}({{ if .RequestType }}req {{ .RequestType }}{{ end }}) error {
	// Call the actual business logic implementation
	// If not implemented, build will fail with: undefined: {{ $.Package | lowercase }}.Impl{{ $.Message.Name }}{{ .Name }}
	return {{ $.Package | lowercase }}.Impl{{ $.Message.Name }}{{ .Name }}(e{{ if .RequestType }}, req{{ end }})
}
{{ end }}

// Getters
{{- range .Message.Fields }}
func (e {{ $.Message.Name }}) {{ .Name | getTitleCase }}() {{ .Type | goType }} {
	return e.{{ .Name | camelCase }}
}
{{- end }}
```

### 5. Generate Business Logic Contracts and Implement Them

For each domain behavior or command:
1. Generate contract stubs in `/business_logic/stubs/`
2. Implement the contracts in `/business_logic/`
3. The build will fail until all required contracts are implemented

**Package Structure and Linking**:
```
/internal/trainings/
  /domain/
    /training/
      training.go          # Generated - imports business_logic package
      repository.go        # Generated interface
  /app/
    /command/
      cancel_training.go   # Generated handler
      
/business_logic/
  /training/
    cancel.go              # Manual implementation - exported functions
  /stubs/
    cancel_stub.go         # Generated reference stub (for documentation)
```

**Build Failure Mechanism**: 
The generated code in `internal/trainings/domain/training/training.go` will import and call functions from the `business_logic/training` package. If those functions are not implemented, Go compilation will fail with standard "undefined: implTrainingCancel" errors. This provides immediate, clear feedback during development.

Example contract stub: `business_logic/stubs/cancel_stub.go`

```go
// Code generated by proto-ddd. DO NOT EDIT.
// This file is a reference stub showing the required contract.
// Implement this function in business_logic/training/cancel.go

package stubs

import (
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/domain/training"
)

// CONTRACT: The following function signature must be implemented
//
// Package: business_logic/training
// Function: ImplTrainingCancel (must be exported)
// Signature: func ImplTrainingCancel(tr *training.Training) error
//
// Description: Implements cancel business logic for the Training entity
//
// Validation rules:
// - Cannot cancel if too close to training time (< 24 hours)
// - Cannot cancel if already canceled
//
// The generated code will import business_logic/training and call this function.
// If not implemented, the build will fail with: undefined: training.ImplTrainingCancel
```
</text>

<old_text line=371>
Example implementation: `business_logic/training/cancel.go`

```go
package training

import (
	"time"
	
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/domain/training"
	"github.com/pkg/errors"
)

// ImplTrainingCancel implements the Cancel method for the Training entity
// This function is called by the generated domain code and must be exported
func ImplTrainingCancel(tr *training.Training) error {
	if tr.IsCanceled() {
		return errors.New("training is already canceled")
	}
	
	// Business rule: can't cancel less than 24h before training
	if time.Until(tr.Time()) < 24*time.Hour {
		return errors.New("can't cancel training less than 24h before it starts")
	}
	
	// Set the canceled flag
	tr.SetCanceled(true)
	
	return nil
}
```

Example implementation: `business_logic/training/cancel.go`

```go
package training

import (
	"time"
	
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/domain/training"
	"github.com/pkg/errors"
)

// Implementation of the Cancel method for the Training entity
func implTrainingCancel(tr *training.Training) error {
	if tr.IsCanceled() {
		return errors.New("training is already canceled")
	}
	
	// Business rule: can't cancel less than 24h before training
	if time.Until(tr.Time()) < 24*time.Hour {
		return errors.New("can't cancel training less than 24h before it starts")
	}
	
	// Set the canceled flag
	tr.SetCanceled(true)
	
	return nil
}
```

### 6. Create Code Generation Tool

**Tool Architecture**: Build a `protoc` plugin called `protoc-gen-ddd` that:
   - Parses proto definitions with custom DDD options
   - Applies templates from the `/templates` directory
   - Generates code in the correct locations within `/internal`
   - Generates business logic contract stubs in `/business_logic/stubs`
   - Preserves manual modifications to business logic files

**Integration with Build Process**:
```bash
# Add to Makefile or build script
protoc \
  --proto_path=proto \
  --ddd_out=internal \
  --ddd_opt=templates=templates \
  --ddd_opt=business_logic=business_logic \
  proto/domain/*.proto \
  proto/application/*.proto
```

**Enforcement Mechanism**:
- Generated domain code imports from `business_logic/training` package
- Calls exported functions like `training.ImplTrainingCancel()`
- If functions don't exist, Go compiler fails with clear errors
- Stub files provide reference documentation for required implementations
- No custom linting needed - standard Go compilation provides enforcement

### 7. Refactor One Domain Area as Proof of Concept

Start with the Training domain:
1. Define all protos
2. Extract business logic
3. Generate code
4. Verify functionality
5. Run tests

### 8. Extend to Other Domain Areas

Once the Training domain is working, apply the same process to:
1. Trainer domain
2. User domain

### 9. Update Documentation and Tests

1. Document the new approach
2. Update tests to work with the new structure
3. Add examples of extending the system

## Success Criteria

- All tests pass after refactoring
- Business logic is clearly separated from infrastructure code
- Adding new features requires minimal boilerplate
- The codebase is more maintainable and follows consistent patterns
- Builds fail with clear error messages when business logic is missing
- Contract-driven approach ensures all required implementations exist

## Timeline

1. **Week 1**: Define proto options and entity schemas
2. **Week 2**: Create basic templates and extraction strategy
3. **Week 3**: Proof of concept with Training domain
4. **Week 4**: Extend to other domains and finalize

## Implementation Checklist

This checklist provides a systematic approach to the refactoring, focusing on end-to-end feature implementation rather than layer-by-layer development.

### Phase 1: Setup and Infrastructure

**Note**: Focus on getting something working to validate the approach, not on perfection.

- [x] Set up directory structure for the new approach
  - [x] Create `/proto` directory and subdirectories
  - [x] Create `/templates` directory and subdirectories
  - [x] Create `/business_logic` directory
  - [x] Create `/business_logic/stubs` directory for contract stubs

- [x] Define custom DDD options
  - [x] Create `ddd_options.proto` with entity, value object, and aggregate root markers
  - [x] Add field-level options (identifier, required, invariant)
  - [x] Add method-level options (validation, authorization)
  - [x] Add application-specific options (dependencies, execution flow)
  - [x] Define structured ExecutionFlow and ExecutionStep messages for type safety
  - [x] Add requires_implementation option for contract enforcement (default to true for all commands)
  - [ ] Design contract interface patterns for business logic implementation

- [ ] Build template processing infrastructure
  - [ ] Create basic template parsing functions
  - [ ] Implement proto parsing with custom option support
  - [ ] Build code generation framework
  - [ ] Add contract verification system with clear error messaging
  - [ ] Create stub generator with documentation comments and validation rules
  - [ ] Implement build-time validation of business logic implementations
  - [ ] Design linking mechanism between generated code and business logic implementations

### Phase 2: Training Domain - Cancel Feature

**Goal**: Create one complete end-to-end example to validate the entire approach.

- [ ] Analyze existing Cancel training implementation
  - [ ] Identify core business rules in the current implementation
  - [ ] Document dependencies and external interactions
  - [ ] Map current code structure to new approach

- [ ] Define proto models for Cancel feature
  - [x] Create core `Training` entity in proto format
  - [x] Define `Cancel` method with appropriate options
  - [ ] Define `CancelTrainingCommand` in application layer

- [ ] Generate business logic contracts and implement them
  - [ ] Generate stub for cancel business logic in `business_logic/stubs/cancel_stub.go`
  - [ ] Implement contract in `business_logic/training/cancel.go`
  - [ ] Verify build fails with clear error messages when implementation is missing
  - [ ] Verify error messages point to exact missing implementations
  - [ ] Ensure all edge cases are preserved
  - [ ] Document contract interface patterns discovered during implementation

- [ ] Create templates for Cancel feature
  - [ ] Entity template with cancel method
  - [ ] Command handler template for cancel operation
  - [ ] Repository interface template

- [ ] Generate code and compare with existing
  - [ ] Generate domain model code
  - [ ] Generate application layer code
  - [ ] Generate adapter code
  - [ ] Verify functionality matches original

### Phase 3: Training Domain - Reschedule Feature

- [ ] Analyze existing Reschedule training implementation
  - [ ] Identify core business rules in the current implementation
  - [ ] Document dependencies and external interactions

- [ ] Define proto models for Reschedule feature
  - [ ] Add `Reschedule` method to `Training` entity
  - [ ] Define `RescheduleTrainingCommand` in application layer

- [ ] Generate business logic contracts and implement them
  - [ ] Generate stub for reschedule business logic in `business_logic/stubs/reschedule_stub.go`
  - [ ] Implement contract in `business_logic/training/reschedule.go`
  - [ ] Verify build fails when implementation is missing
  - [ ] Validate contract interface consistency with the Cancel feature
  - [ ] Note any emerging patterns for potential future default implementations

- [ ] Update templates if needed based on learnings
  - [ ] Refine entity templates
  - [ ] Refine command handler templates
  - [ ] Refine contract stub templates with better documentation
  - [ ] Identify any common patterns that could be templated in future iterations

- [ ] Generate code and compare with existing
  - [ ] Verify functionality matches original

### Phase 4: Training Domain - Query Features

- [ ] Analyze existing query implementations
  - [ ] Map current query structure to new approach

- [ ] Define proto models for queries
  - [ ] Define `TrainingsForUserQuery` and related types

- [ ] Create templates for query operations
  - [ ] Query handler templates
  - [ ] Read model templates

- [ ] Generate code and compare with existing
  - [ ] Verify query functionality matches original

### Phase 5: Optimize and Refine Based on Patterns

**Stopping point evaluation**: After completing the Training domain, evaluate whether:
- The approach is technically feasible
- The ergonomics feel right for developers
- The separation of concerns is actually clearer
- It's worth continuing to other domains

- [ ] Review patterns identified in Training domain
  - [ ] Update templates to better accommodate common patterns
  - [ ] Refactor proto definitions if needed
  - [ ] Consider adding default implementations for common patterns
  - [ ] Document contract interfaces that could benefit from standardization
  - [ ] Document lessons learned and whether to proceed

- [ ] Improve code generation process
  - [ ] Add support for incremental regeneration
  - [ ] Implement protection for manual modifications

### Phase 6: Extend to Trainer Domain

- [ ] Apply same process to Trainer domain
  - [ ] Define proto models
  - [ ] Extract business logic
  - [ ] Generate and verify code

### Phase 7: Extend to User Domain

- [ ] Apply same process to User domain
  - [ ] Define proto models
  - [ ] Extract business logic
  - [ ] Generate and verify code

### Phase 8: Documentation and Tests

- [ ] Update documentation
  - [ ] Document the new approach
  - [ ] Create guides for adding new features

- [ ] Ensure all tests pass
  - [ ] Update tests if needed
  - [ ] Add tests for the template system

### Phase 9: Finalization

- [ ] Clean up any remaining technical debt
- [ ] Perform final verification of all features
- [ ] Complete final documentation
- [ ] Document patterns for potential default implementations in future versions
- [ ] Create guide for adding new features using the contract-driven approach