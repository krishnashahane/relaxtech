package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/krishna/relaxtech/internal/client"
	"github.com/krishna/relaxtech/internal/output"
)

// travelCmd groups travel and jetlag management subcommands.
var travelCmd = &cobra.Command{Use: "travel", Short: "Travel / jetlag endpoints"}

// travelTripsCmd lists all saved trips.
var travelTripsCmd = &cobra.Command{Use: "trips", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	tripList, err := cl.Travel().Trips(context.Background())
	if err != nil {
		return err
	}
	return output.Render(output.DisplayMode(viper.GetString("output")), []string{"trips"}, []map[string]any{{"trips": tripList}})
}}

// travelCreateTripCmd creates a new trip entry.
var travelCreateTripCmd = &cobra.Command{Use: "create-trip", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	tripData := map[string]any{}
	if dest := viper.GetString("destination"); dest != "" {
		tripData["destination"] = dest
	}
	if begin := viper.GetString("start-date"); begin != "" {
		tripData["startDate"] = begin
	}
	if finish := viper.GetString("end-date"); finish != "" {
		tripData["endDate"] = finish
	}
	if tz := viper.GetString("timezone"); tz != "" {
		tripData["timezone"] = tz
	}
	if len(tripData) == 0 {
		return fmt.Errorf("provide at least --destination or --start-date/--end-date")
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.Travel().CreateTrip(context.Background(), tripData)
}}

// travelDeleteTripCmd removes a trip by its ID.
var travelDeleteTripCmd = &cobra.Command{Use: "delete-trip", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	tripID := viper.GetString("trip")
	if tripID == "" {
		return fmt.Errorf("--trip required")
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.Travel().DeleteTrip(context.Background(), tripID)
}}

// travelPlansCmd lists plans associated with a trip.
var travelPlansCmd = &cobra.Command{Use: "plans", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	tripID := viper.GetString("trip")
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	planList, err := cl.Travel().Plans(context.Background(), tripID)
	if err != nil {
		return err
	}
	return output.Render(output.DisplayMode(viper.GetString("output")), []string{"plans"}, []map[string]any{{"plans": planList}})
}}

// travelCreatePlanCmd creates a new plan within a trip.
var travelCreatePlanCmd = &cobra.Command{Use: "create-plan", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	tripID := viper.GetString("trip")
	if tripID == "" {
		return fmt.Errorf("--trip required")
	}
	planData := map[string]any{}
	if label := viper.GetString("name"); label != "" {
		planData["name"] = label
	}
	if when := viper.GetString("date"); when != "" {
		planData["date"] = when
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.Travel().CreatePlan(context.Background(), tripID, planData)
}}

// travelUpdatePlanCmd modifies an existing travel plan.
var travelUpdatePlanCmd = &cobra.Command{Use: "update-plan", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	planID := viper.GetString("plan")
	if planID == "" {
		return fmt.Errorf("--plan required")
	}
	changes := map[string]any{}
	if label := viper.GetString("name"); label != "" {
		changes["name"] = label
	}
	if when := viper.GetString("date"); when != "" {
		changes["date"] = when
	}
	if len(changes) == 0 {
		return fmt.Errorf("no fields to update")
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.Travel().UpdatePlan(context.Background(), planID, changes)
}}

// travelTasksCmd lists tasks associated with a travel plan.
var travelTasksCmd = &cobra.Command{Use: "tasks", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	planID := viper.GetString("plan")
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	taskList, err := cl.Travel().PlanTasks(context.Background(), planID)
	if err != nil {
		return err
	}
	return output.Render(output.DisplayMode(viper.GetString("output")), []string{"tasks"}, []map[string]any{{"tasks": taskList}})
}}

// travelAirportCmd searches for airports by keyword.
var travelAirportCmd = &cobra.Command{Use: "airport-search", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	searchTerm := viper.GetString("query")
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	matches, err := cl.Travel().AirportSearch(context.Background(), searchTerm)
	if err != nil {
		return err
	}
	return output.Render(output.DisplayMode(viper.GetString("output")), []string{"airports"}, []map[string]any{{"airports": matches}})
}}

// travelFlightCmd retrieves the status of a specific flight.
var travelFlightCmd = &cobra.Command{Use: "flight-status", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	flightNum := viper.GetString("flight")
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	flightInfo, err := cl.Travel().FlightStatus(context.Background(), flightNum)
	if err != nil {
		return err
	}
	return output.Render(output.DisplayMode(viper.GetString("output")), []string{"flight"}, []map[string]any{{"flight": flightInfo}})
}}

func init() {
	travelPlansCmd.Flags().String("trip", "", "trip id")
	travelTasksCmd.Flags().String("plan", "", "plan id")
	travelAirportCmd.Flags().String("query", "", "airport query")
	travelFlightCmd.Flags().String("flight", "", "flight number")
	travelCreateTripCmd.Flags().String("destination", "", "destination")
	travelCreateTripCmd.Flags().String("start-date", "", "start date")
	travelCreateTripCmd.Flags().String("end-date", "", "end date")
	travelCreateTripCmd.Flags().String("timezone", "", "timezone")
	travelDeleteTripCmd.Flags().String("trip", "", "trip id")
	travelCreatePlanCmd.Flags().String("trip", "", "trip id")
	travelCreatePlanCmd.Flags().String("name", "", "plan name")
	travelCreatePlanCmd.Flags().String("date", "", "plan date")
	travelUpdatePlanCmd.Flags().String("plan", "", "plan id")
	travelUpdatePlanCmd.Flags().String("name", "", "plan name")
	travelUpdatePlanCmd.Flags().String("date", "", "plan date")

	viper.BindPFlag("trip", travelPlansCmd.Flags().Lookup("trip"))
	viper.BindPFlag("plan", travelTasksCmd.Flags().Lookup("plan"))
	viper.BindPFlag("query", travelAirportCmd.Flags().Lookup("query"))
	viper.BindPFlag("flight", travelFlightCmd.Flags().Lookup("flight"))
	viper.BindPFlag("destination", travelCreateTripCmd.Flags().Lookup("destination"))
	viper.BindPFlag("start-date", travelCreateTripCmd.Flags().Lookup("start-date"))
	viper.BindPFlag("end-date", travelCreateTripCmd.Flags().Lookup("end-date"))
	viper.BindPFlag("timezone", travelCreateTripCmd.Flags().Lookup("timezone"))
	viper.BindPFlag("trip", travelDeleteTripCmd.Flags().Lookup("trip"))
	viper.BindPFlag("trip", travelCreatePlanCmd.Flags().Lookup("trip"))
	viper.BindPFlag("name", travelCreatePlanCmd.Flags().Lookup("name"))
	viper.BindPFlag("date", travelCreatePlanCmd.Flags().Lookup("date"))
	viper.BindPFlag("plan", travelUpdatePlanCmd.Flags().Lookup("plan"))
	viper.BindPFlag("name", travelUpdatePlanCmd.Flags().Lookup("name"))
	viper.BindPFlag("date", travelUpdatePlanCmd.Flags().Lookup("date"))

	travelCmd.AddCommand(
		travelTripsCmd,
		travelCreateTripCmd,
		travelDeleteTripCmd,
		travelPlansCmd,
		travelCreatePlanCmd,
		travelUpdatePlanCmd,
		travelTasksCmd,
		travelAirportCmd,
		travelFlightCmd,
	)
}
